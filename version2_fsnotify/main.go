package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	shell "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
)

var FileCidMap map[string]string

type Watch struct {
	watch *fsnotify.Watcher
}

//watch dir
func (w *Watch) watchDir(dir string) {
	FileCidMap = make(map[string]string)

	sh := shell.NewShell("localhost:5001")
	log.Info("ipfs start")

	c := context.Background()

	//use Walk function to range son dir
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("watch : ", path)
		}
		return nil
	})
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("create file : ", ev.Name)
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name)
							fmt.Println("add watch : ", ev.Name)

							//ipfs add dir
							opt := func(builder *shell.RequestBuilder) error {
								builder.Option("p", true)
								return nil
							}

							err = sh.FilesMkdir(c, ev.Name, opt)
							if err != nil {
								log.Error("ipfs make dir failed: ", err)
							}
							log.Info("ipfs make dir success: ", ev.Name)
						} else {
							//ipfs add file
							data, err := ioutil.ReadFile(ev.Name)
							if err != nil {
								log.Error("Read File error: %s", err)
								continue
							}
							cid, err := sh.Add(strings.NewReader(string(data)))
							if err != nil {
								log.Error("ipfs add failed: ", err)
							}
							log.Info("ipfs add file:"+ev.Name+" cid is: ", cid)
							FileCidMap[ev.Name] = cid
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("update file : ", ev.Name)
						data, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							log.Error("Read File error: %s", err)
							continue
						}
						// opt := func(builder *shell.RequestBuilder) error {
						// 	builder.Option("p", true)
						// 	return nil
						// }
						// err = sh.FilesWrite(c, ev.Name, strings.NewReader(string(data)), opt) //TODO:如何获取命令行输入的内容
						// if err != nil {
						// 	log.Error("ipfs write file failed: ", ev.Name)
						// 	continue
						// }
						cid, err := sh.Add(strings.NewReader(string(data)))
						if err != nil {
							log.Error("ipfs modify failed: ", err)
							continue
						}
						log.Info("ipfs modify file:"+ev.Name+" cid is: ", cid)
						log.Info("ipfs write file success:" + ev.Name)
						//FileCidMap[ev.Name] = cid
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("删除文件 : ", ev.Name)

						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() { // this name is dir
							w.watch.Remove(ev.Name)
							fmt.Println("删除监控 : ", ev.Name)
							err = sh.FilesRm(c, ev.Name, false) //TODO:参数true代表删除目录是否递归？
							if err != nil {
								log.Error("ipfs rm dir failed: ", err)
								continue
							}
							log.Info("ipfs rm dir success: ", ev.Name)
						} else {
							err = sh.FilesRm(c, ev.Name, false)
							if err != nil {
								log.Error("ipfs rm file failed: ", err)
								continue
							}
							log.Info("ipfs rm success: ", ev.Name)
						}

						delete(FileCidMap, ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("rename file : ", ev.Name)
						w.watch.Remove(ev.Name)
						//ipfs rm
						opt := func(builder *shell.RequestBuilder) error {
							builder.Option("p", true)
							return nil
						}
						fi, err := sh.FilesLs(c, ev.Name, opt)
						if len(fi) == 0 {
							log.Error("the reason why file not exist：", err)
							continue
						}

						err = sh.FilesRm(c, ev.Name, false)
						if err != nil {
							log.Error("ipfs rm failed: ", err)
							//continue
						}
						log.Info("ipfs rm success: ", ev.Name)
						delete(FileCidMap, ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("chmod permit : ", ev.Name)

					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()
}

func main() {
	//0，init
	filePath := "/home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test"
	fmt.Println("ipfs watch dir is: ", filePath)

	//FileCidMap := make(map[string]string)

	log.Info("start")

	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	w.watchDir(filePath)
	//loop
	select {}
}

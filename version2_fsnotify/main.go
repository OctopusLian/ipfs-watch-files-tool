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

var FileCidMap map[string]string //var声明，不占用内存

type Watch struct {
	watch *fsnotify.Watcher
}

//监控目录
func (w *Watch) watchDir(dir string) {
	//FileCidMap = make(map[string]string)

	sh := shell.NewShell("localhost:5001")
	log.Info("ipfs start")

	c := context.Background()

	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控 : ", path)
		}
		return nil
	})
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("创建文件 : ", ev.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name)
							fmt.Println("添加监控 : ", ev.Name)

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
							//FileCidMap[ev.Name] = cid
						}

					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						fmt.Println("写入文件 : ", ev.Name)
						data, err := ioutil.ReadFile(ev.Name)
						if err != nil {
							log.Error("Read File error: %s", err)
							continue
						}
						// cid, err := sh.Add(strings.NewReader(string(data)))
						// if err != nil {
						// 	log.Error("ipfs add failed: ", err)
						// }
						opt := func(builder *shell.RequestBuilder) error {
							builder.Option("p", true)
							return nil
						}
						err = sh.FilesWrite(c, ev.Name, strings.NewReader(string(data)), opt) //TODO:如何获取命令行输入的内容
						if err != nil {
							log.Error("ipfs write file failed: ", ev.Name)
							continue
						}
						log.Info("ipfs write file success:" + ev.Name)
						//FileCidMap[ev.Name] = cid
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("删除文件 : ", ev.Name)
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() { // this name is dir
							w.watch.Remove(ev.Name)
							fmt.Println("删除监控 : ", ev.Name)
							err = sh.FilesRm(c, ev.Name, true) //TODO:参数true代表删除目录是否递归？
							if err != nil {
								log.Error("ipfs rm dir failed: ", err)
								continue
							}
							log.Info("ipfs rm dir success: ", ev.Name)
						} else {
							//是文件，就不需要移除监控，直接ipfs调用删除接口
							err = sh.FilesRm(c, ev.Name, false)
							if err != nil {
								log.Error("ipfs rm file failed: ", err)
								continue
							}
							log.Info("ipfs rm success: ", ev.Name)
						}

						//delete(FileCidMap, ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						fmt.Println("重命名文件 : ", ev.Name)
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						w.watch.Remove(ev.Name)
						//ipfs rm
						err := sh.FilesRm(c, ev.Name, true)
						if err != nil {
							log.Error("ipfs rm failed: ", err)
						}
						log.Info("ipfs rm success: ", ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						fmt.Println("修改权限 : ", ev.Name)

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
	//0，准备工作
	filePath := "/home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test"
	fmt.Println("ipfs watch dir is: ", filePath)

	//FileCidMap := make(map[string]string)

	log.Info("start")

	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	w.watchDir(filePath)
	//循环
	select {}
}

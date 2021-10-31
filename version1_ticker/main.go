package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
)

func main() {
	//0，准备工作
	filePath := "/home/neo/Code/go/src/github.com/OctopusLian/ipfs-watch-files-tool/test"
	fmt.Println("ipfs watch dir is: ", filePath)

	FileCidMap := make(map[string]string)

	sh := shell.NewShell("localhost:5001")
	log.Info("start")

	ticker := time.NewTicker(10 * time.Second) //10s启动一次
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done: //监听停止通道
			fmt.Println("Done!")
			return
		case t := <-ticker.C: //开始做ipfs监控的任务
			log.Info("Current time: ", t)

			//1，读取监控路径下的文件，返回是一个列表
			files, err := ioutil.ReadDir(filePath)
			if err != nil {
				log.Error("ReadDir failed: ", err)
				return
			}

			//2，将读取到的文件，存入ipfs中
			for _, file := range files {
				if _, ok := FileCidMap[file.Name()]; !ok {
					fileNamePath := filePath + "/" + file.Name()
					fmt.Println("watch file name is: ", fileNamePath)
					//map里面没有
					data, err := ioutil.ReadFile(fileNamePath)
					if err != nil {
						log.Error("Read File error: %s", err)
						continue
					}
					cid, err := sh.Add(strings.NewReader(string(data)))
					if err != nil {
						log.Error(os.Stderr, "error: %s", err)
						os.Exit(1)
					}
					log.Info(fileNamePath+" cid is: ", cid)

					FileCidMap[fileNamePath] = cid
				}
			}

		}
	}
}

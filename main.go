package main

import (
	"time"
	"log"
	"./config"
	"os"
	"fmt"
	"math/rand"
)

/*
1.第一次初始化配置文件内容
2.可以临时修改文件内容
3.删除文件不影响程序正常运行
下个版本添加删除配置文件自动恢复到最后一个版本
*/

func main() {
	log.Println("main:", config.Conf)
	ticker := time.NewTicker(time.Second * 1)
	ticker1 := time.NewTicker(time.Second * 2)
	ticker2 := time.NewTicker(time.Second * 20)
	go func() {
		for _ = range ticker.C {
			log.Println("main:", config.Conf)
		}
	}()
	for {

		select {
		case <-ticker1.C:
			if file, err := os.OpenFile(config.FileName, os.O_APPEND, 0666); err == nil {
				file.WriteString(fmt.Sprintf("%d = %d   \n", rand.Int(), rand.Int()))
				file.Close()
			} else {
				log.Println("main（）打开文件不存在",err)
				continue
			}
		case <-ticker2.C:
			if removeErr := os.Remove(config.FileName); removeErr == nil {
				log.Println("删除文件成功")
			} else {
				log.Println("删除文件失败")
			}
		}
	}
}

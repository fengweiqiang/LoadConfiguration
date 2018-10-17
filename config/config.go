package config

import (
	"sync"
	"time"
	"os"
	"bufio"
	"fmt"
	"log"
	"io"
	"strings"
)

/*
LoadConfiguration.conf

aaa=222
bbb=333
asd

 */
var Conf map[string]string
var lock sync.Mutex

var FileName = "LoadConfiguration.conf"
//第一次启动初始化文件内容
var Content=""
func init() {
	ticker := time.NewTicker(time.Second * 5)
	var lastTime int64
	if _, err := os.Stat(FileName); os.IsNotExist(err) {
		if file, err := os.Create(FileName); err == nil {
			w := bufio.NewWriter(file)
			fmt.Fprintln(w, Content)
			w.Flush()
			defer file.Close()
			log.Println(FileName, "创建写入成功")
		} else {
			log.Println(FileName, "创建写入失败")
		}
	} else {
		log.Println(FileName, "已存在")
	}
	go func() {
		for _ = range ticker.C {
			file, err := os.OpenFile(FileName, os.O_CREATE|os.O_RDONLY, 0666)
			if err != nil {
				log.Println("OpenFile()打开文件失败，重新创建文件--->")
				os.Create(FileName)
				return
			}
			fileInfo, err := file.Stat()
			file.Close()
			if err != nil {
				fmt.Printf("stat %s failed,err:%v\n", file, err)
				return
			}
			curModifyTime := fileInfo.ModTime().Unix()
			log.Println("Conf:", Conf, len(Conf))
			log.Printf("本次获取修改时间%v --- 上次获取修改时间%v\n", curModifyTime, lastTime)
			//判断文件的修改时间是否大于最后一次修改时间
			if curModifyTime > lastTime {
				m, err := parse()
				if err != nil {
					fmt.Println("parse failed,err:", err)
					return
				}
				lock.Lock()
				Conf = m
				//log.Println(CallBack)
				lock.Unlock()
				lastTime = curModifyTime
			}
		}
	}()
}

func parse() (m map[string]string, err error) {
	// 读文件并或将文件中的数据以k/v的形式存储到map中
	m = make(map[string]string, 20)
	file, err := os.Open(FileName)
	defer file.Close()
	if err != nil {
		log.Println("parse()打开文件失败，重新创建文件--->")
		os.Create(FileName)
		return
	}
	var lineNo int
	reader := bufio.NewReader(file)
	for {
		// 一行行的读文件
		line, errRet := reader.ReadString('\n')
		log.Println("line:", line, errRet)
		if len(line) == 0 && errRet == io.EOF {
			return
		}
		lineNo++
		line = strings.TrimSpace(line) // 取出空格
		//log.Println("lineaaa:",line)
		if len(line) == 0 || line[0] == '\n' || line[0] == '+' || line[0] == ';' {
			// 当前行为空行或者是注释行等
			continue
		}
		arr := strings.Split(line, "=") // 通过=进行切割取出k/v结构
		if len(arr) <= 1 {
			log.Printf("invalid config,line:%d\n", lineNo)
			continue
		}
		key := strings.TrimSpace(arr[0])
		if len(key) == 0 {
			log.Printf("invalid config,line:%d\n", lineNo)
			continue
		}
		if len(arr) == 1 {
			m[key] = ""
			continue
		}
		value := strings.TrimSpace(arr[1])
		m[key] = value
	}
	return
}



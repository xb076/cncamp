package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	const FILE_PATH_CONF string = "config.conf"
	const FILE_PATH_LOG string = "log.log"
	myLog.Start(FILE_PATH_LOG)

	f := new(FileOps)
	f.Open(FILE_PATH_CONF, 0)
	conf := f.FileConfig2Server(FILE_PATH_CONF)
	f.Close()

	fmt.Println("Server start ... Enter Q to exit")
	myLog.LogInfo("Server start ...")
	svr := new(Server)
	go svr.Start(conf.server_ip, conf.server_port)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入Q就退出
			break
		}
	}

	fmt.Println("Q entered, server exit ... ")
	myLog.LogInfo("Q entered, server exit ... ")
	myLog.Close()
}

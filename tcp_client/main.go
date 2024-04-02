package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TCP 客户端
func main() {
	const FILE_PATH_CONF string = "config.conf"
	const FILE_PATH_LOG string = "log.log"
	myLog.Start(FILE_PATH_LOG)

	f := new(FileOps)
	f.Open(FILE_PATH_CONF, 0)
	conf := f.FileConfig2Client(FILE_PATH_CONF)
	f.Close()

	cc := new(ClientCore)
	cc.Start(conf.server_ip, conf.server_port, conf.client_thread)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" { // 如果输入q就退出
			break
		}
	}

	cc.Stop()
	fmt.Println("Client exit")
	myLog.Close()
}

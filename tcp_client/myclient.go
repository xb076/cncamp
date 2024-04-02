package main

import (
	"fmt"
	"net"
	"time"
)

type ClientCore struct {
	server_ip    string
	server_port  string
	client_num   int
	client_slice []*Client
}

func (cc *ClientCore) Start(ip string, port string, client_num int) {
	cc.server_ip = ip
	cc.server_port = port
	cc.client_num = client_num
	myLog.LogInfo("Client start")
	for i := 0; i < cc.client_num; i++ {
		client := new(Client)
		client.num_sn = i
		client.server_ip = cc.server_ip
		client.server_port = cc.server_port
		client.milli = 1000
		client.msg_idx = 0
		cc.client_slice = append(cc.client_slice, client)
		go client.Start()
	}
}

func (cc *ClientCore) Stop() {
	for i := 0; i < cc.client_num; i++ {
		client := cc.client_slice[i]
		client.conn.Close()
	}
	myLog.LogInfo("Client closed")
}

type Client struct {
	num_sn      int
	server_ip   string
	server_port string
	conn        net.Conn
	milli       int64
	buffer      [128]byte
	msg_idx     int
}

func (client *Client) Start() {
	var err error
	address := client.server_ip + ":" + client.server_port
	client.conn, err = net.Dial("tcp", address)
	if nil != err {
		msg := fmt.Sprintf("Client error: <%s>", err.Error())
		fmt.Println(msg)
		myLog.LogError(msg)
		return
	}
	ticker := time.NewTicker(time.Duration(client.milli) * time.Millisecond)
	for _ = range ticker.C {
		msg := fmt.Sprintf("client message: client id: %d, msg id: %d", client.num_sn, client.msg_idx)
		client.msg_idx++
		_, err = client.conn.Write([]byte(msg))
		if nil != err {
			msg := fmt.Sprintf("Client error: <%s>", err.Error())
			fmt.Println(msg)
			myLog.LogError(msg)
			break
		}
		n, err := client.conn.Read(client.buffer[:])
		if err != nil {
			msg := fmt.Sprintf("Client error: <%s>", err.Error())
			fmt.Println(msg)
			myLog.LogError(msg)
			break
		}
		msg = fmt.Sprintf("Client#<%d> Receive: <%s>", client.num_sn, (client.buffer[:n]))
		fmt.Println(msg)
		myLog.LogInfo(msg)
		emptyBytes(client.buffer[:], n)
	}

}

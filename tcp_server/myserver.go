package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	address      string
	thread_index int
	thread_num   int
	thread_slice []*ServerThread
}

func (svr *Server) GetThread() *ServerThread {
	var svr_t *ServerThread
	if svr.thread_num > 0 {
		for i := 0; i < svr.thread_num; i++ {
			if svr.thread_slice[i].status == 0 {
				svr_t = svr.thread_slice[i]
				svr.thread_index = i
				return svr_t
			}
		}
	}
	svr_t = new(ServerThread)
	svr_t.num_sn = svr.thread_num
	svr_t.status = 0
	svr_t.sig = make(chan bool, 1)
	svr.thread_slice = append(svr.thread_slice, svr_t)
	go svr_t.Start()
	svr.thread_index = svr.thread_num
	svr.thread_num += 1

	return svr_t
}

func (svr *Server) Start(ip string, port string) {
	svr.address = ip + ":" + port
	svr.thread_index = 0

	listen, err := net.Listen("tcp", svr.address)
	if err != nil {
		fmt.Println("Listen() failed, err: ", err)
		myLog.LogError("Listen() failed, err: " + err.Error())
		return
	}
	fmt.Println("Server listen: ", svr.address)
	myLog.LogInfo("Server listen: " + svr.address)
	svr_t := svr.GetThread()
	for {
		svr_t.conn, err = listen.Accept() // 监听客户端的连接请求
		if err != nil {
			myLog.LogError("Accept() error: " + err.Error())
			continue
		}
		svr_t.sig <- true
		svr_t.status = 1
		svr_t.remoteAdd = svr_t.conn.RemoteAddr().String()
		svr_t.b_exit = false
		msg := fmt.Sprintf("Received from remote <%s>", svr_t.remoteAdd)
		fmt.Println(msg)
		myLog.LogInfo(msg)
		svr_t = svr.GetThread()
	}

}

type ServerThread struct {
	num_sn    int
	buffer    [1024]byte
	conn      net.Conn
	status    int //0 idle, 1 busy
	sig       chan bool
	b_exit    bool
	msg_idx   int
	remoteAdd string
}

func (svr_t *ServerThread) Start() {
	myLog.LogInfo(fmt.Sprintf("Thread#%d start", svr_t.num_sn))
	svr_t.msg_idx = 0
	for {
		if svr_t.status == 0 {
			msg := fmt.Sprintf("Thread#%d waiting", svr_t.num_sn)
			fmt.Println(msg)
			myLog.LogInfo(msg)
			<-svr_t.sig
		}
		svr_t.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := svr_t.conn.Read(svr_t.buffer[:])
		if nil != err {
			str := err.Error()
			if strings.Contains(str, "i/o timeout") {
				if svr_t.b_exit {
					break
				}
				continue
			} else if strings.Contains(str, "EOF") {
				svr_t.conn.Close()
				svr_t.status = 0
			}
			msg := fmt.Sprintf("read from client failed, err: <%s>", err)
			fmt.Println(msg)
			myLog.LogError(msg)
			continue
		}
		recvStr := string(svr_t.buffer[:n])
		msg := fmt.Sprintf("Thread#<%d> Recieve: <%s>", svr_t.num_sn, recvStr)
		fmt.Println(msg)
		myLog.LogInfo(msg)
		msg = fmt.Sprintf("Server msg: thread id: %d, msg id: %d", svr_t.num_sn, svr_t.msg_idx)
		svr_t.conn.Write([]byte(msg))
		myLog.LogInfo(msg)
		svr_t.msg_idx++
	}
}

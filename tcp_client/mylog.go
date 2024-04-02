package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var myLog LogFile

type LogFile struct {
	filepath string
	fileptr  *os.File
	ch_msg   chan []byte
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func init() {
	myLog.ctx, myLog.cancel = context.WithCancel(context.Background())
	myLog.ch_msg = make(chan []byte)
	//myLog.ch_sig = make(chan bool)
	myLog.wg = sync.WaitGroup{}
}

func formatTimeString(bytes []byte) string {
	msg := string(bytes)
	end := strings.Index(msg, "]")
	n := strings.Index(msg, ".")
	result := string(bytes[0 : n+5])
	result += string(bytes[end:])
	return result
}

func (l *LogFile) Start(path string) {
	var err error
	l.filepath = path
	l.fileptr, err = os.OpenFile(l.filepath, os.O_APPEND, 0666)
	if nil != err {
		fmt.Printf("Log file open error: <%s>\n", err)
		return
	}

	go l.Logging()
	fmt.Printf("Log file <%s> open success\n", l.filepath)
}

func (l *LogFile) LogInfo(str string) {
	l.wg.Add(1)
	msg := "[" + time.Now().String() + "]" + "<INFO>" + str + "\n"
	l.ch_msg <- []byte(msg)
	//fmt.Printf("Log msg: %s", msg)
}

func (l *LogFile) LogWarn(str string) {
	l.wg.Add(1)
	msg := "[" + time.Now().String() + "]" + "<WARN>" + str + "\n"
	l.ch_msg <- []byte(msg)
	//fmt.Printf("Log msg: %s", msg)
}

func (l *LogFile) LogError(str string) {
	l.wg.Add(1)
	msg := "[" + time.Now().String() + "]" + "<ERROR>" + str + "\n"
	l.ch_msg <- []byte(msg)
	//fmt.Printf("Log msg: %s", msg)
}

func (l *LogFile) Logging() {
	var b_exit = false
	for {
		select {
		case data := <-l.ch_msg:
			//fmt.Printf("Log data to write: <%s>\n", data)
			l.Write([]byte(formatTimeString(data)))
		case <-l.ctx.Done():
			b_exit = true
		}
		if b_exit {
			l.wg.Wait()
			fmt.Println("Logging thread exit")
			break
		}
	}
}

func (l *LogFile) Write(data []byte) {
	n, err := l.fileptr.Write([]byte(data))
	if nil != err {
		fmt.Printf("Log file written error: <%d> bytes <%s>\n", n, err)
	}
	l.wg.Done()
}

func (l *LogFile) Close() {
	l.wg.Wait()
	l.cancel()
	l.fileptr.Close()
}

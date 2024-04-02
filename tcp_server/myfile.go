package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const BUFF_SIZE = 6

func emptyBytes(b []byte, length int) {
	if length < 1 {
		length = len(b)
	}

	for i := 0; i < length; i++ {
		b[i] = 0
	}
}

func trimString(s string) string {
	b := []byte(s)
	r := []byte{}
	for i := 0; i < len(b); i++ {
		if b[i] > 32 && b[i] < 127 {
			r = append(r, b[i])
		}
	}
	return string(r)
}

type FileOps struct {
	filepath   string
	fileptr    *os.File
	filesize   int64
	pos_seek   int64
	line_begin int64
	line_end   int64
	pos_offset int64
}

func (f *FileOps) Open(path string, mode int) (err error) {
	f.filepath = path
	if mode == 0 {
		f.fileptr, err = os.OpenFile(f.filepath, os.O_RDONLY, 0666)
	} else if mode == 1 {
		f.fileptr, err = os.OpenFile(f.filepath, os.O_WRONLY, 0666)
	}

	if nil != err {
		fmt.Println(err)
	} else {
		msg := fmt.Sprintf("file '%s' open success", f.filepath)
		fmt.Println(msg)
		//myLog.LogInfo(msg)
	}

	fileInfo, _ := os.Stat(f.filepath)
	f.filesize = fileInfo.Size()

	f.Reset()

	return err
}

func (f *FileOps) Close() {
	if nil != f.fileptr {
		f.fileptr.Close()
		msg := fmt.Sprintf("file '%s' closed", f.filepath)
		fmt.Println(msg)
		//myLog.LogInfo(msg)
	}
}

func (f *FileOps) Reset() (err error) {
	f.pos_seek = 0
	f.line_begin = 0
	f.line_end = 0
	f.pos_offset = 0

	f.pos_offset, err = f.fileptr.Seek(0, 0)

	return err
}

func (f *FileOps) ReadLine() string {
	var line_read string
	var line_length int64
	var err error
	var read_char byte
	var i int64 = 0
	var n int = 0
	var gotLineEnd bool

	buffer := [BUFF_SIZE]byte{0}
	//f.pos_offset = f.line_begin
	emptyBytes(buffer[:], 0)
	line_read = ""
	gotLineEnd = false
	for i = f.line_begin; i < f.filesize; i++ {
		//fmt.Printf("line_begin#1: %d\n", f.line_begin)
		pos := f.line_begin
		_, err = f.fileptr.Seek(pos, 0)
		if nil != err {
			log.Fatal(err)
			break
		}
		n, err = f.fileptr.Read(buffer[:])
		if nil != err {
			log.Fatal(err)
			break
		}
		//fmt.Printf("Buffer read %d bytes: %s\n", n, buffer[0:n])

		f.line_end = f.line_begin
		for j := 0; j < n; j++ {
			read_char = buffer[j]
			f.line_end += 1
			if read_char == 13 || read_char == 10 {
				f.line_end -= 1
				line_length = f.line_end - f.line_begin
				line_read += string(buffer[0:line_length])
				//fmt.Println("read line over")
				//fmt.Printf("read line %d bytes: %s\n", line_length, line_read)
				//fmt.Println("line end position: ", f.line_end)
				if read_char == 13 {
					f.line_begin = f.line_end + 2
				} else if read_char == 10 {
					f.line_begin = f.line_end + 1
				}
				//fmt.Printf("line_begin#2: %d\n", f.line_begin)
				gotLineEnd = true
				break
			} else if j+1 == n {
				line_length = f.line_end - f.line_begin
				line_read += string(buffer[0:line_length])
				f.line_begin = f.line_end
				break
			}
		}
		i = f.line_begin
		emptyBytes(buffer[:], n)

		if gotLineEnd || f.line_end == f.filesize {
			break
		}

	}
	//fmt.Printf("Readline %d bytes: %s\n", len(line_read), line_read)
	return line_read
}

func (f *FileOps) WriteLine(bytes []byte) (n int) {
	var err error
	var n_write int
	//_, err = f.fileptr.Seek(0, 2)
	if nil != err {
		fmt.Println("Error seek() in WriteLine(): ", err)
	}
	n_write, err = f.fileptr.Write(bytes)
	if nil != err {
		fmt.Println("Error write() in WriteLine(): ", err)
	}
	return n_write
}

func (f *FileOps) FileConfig2Server(path string) *ServerConfig {
	var file_line string
	var ip string
	var port string
	conf := new(ServerConfig)
	conf.server_ip = ""
	conf.server_port = ""
	file_line = ""
	for {
		file_line = f.ReadLine()
		//fmt.Printf("ReadLine() by ServerConfig %d bytes: %s\n", len(file_line), file_line)
		if len(file_line) < 1 {
			break
		}
		if strings.Contains(file_line, "IP") {
			if i := strings.Index(file_line, "="); i > 0 {
				line_bytes := []byte(file_line)
				ip = string(line_bytes[i+1:])
				conf.server_ip = trimString(ip)
			}
		}
		if strings.Contains(file_line, "PORT") {
			if i := strings.Index(file_line, "="); i > 0 {
				line_bytes := []byte(file_line)
				port = string(line_bytes[i+1:])
				conf.server_port = trimString(port)
			}
		}

	}
	fmt.Println("Server ip is:", conf.server_ip)
	fmt.Println("Server ip port:", conf.server_port)
	return conf
}

type ServerConfig struct {
	server_ip     string
	server_port   string
	server_thread int
}


/*
生产线程生产10秒后阻塞
消费线程3秒后开始消费
*/

package main

import(
	"fmt"
	"time"
)

func main(){

	buff := make(chan int, 10)
	sig_stopWrite := make(chan bool, 1)
	sig_stopRead := make(chan bool, 1)

	go func(){
		fmt.Println("production process started")
		var ch chan<- int
		ch = buff

		//time.Sleep(3 * time.Second)
		num := 0
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <- sig_stopWrite:
				fmt.Println("production process channel done()")
				return
			default:
				num++
				ch <- num
				fmt.Println("Produce: ", num)
			}
		}
	}()

	go func(){
		fmt.Println("consumption process started")
		var ch <-chan int
		ch = buff

		time.Sleep(12 * time.Second) 
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <- sig_stopRead:
				fmt.Println("consumption process channel done()")
				return
			default:
				fmt.Println("Consume: ", <- ch)
			}
		}
	}()

	defer func(){
		time.Sleep(3 * time.Second)
		close(buff)
		fmt.Println("main process exit")
	}()

	time.Sleep(20 * time.Second)
	close(sig_stopWrite)
	close(sig_stopRead)
	
}



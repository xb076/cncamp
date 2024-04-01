

package main

import(
	"fmt"
	"time"
	"sync"
)

type DataPool struct{
	buff chan int
	size int
	length int
	cond *sync.Cond
	wg sync.WaitGroup
}

type Producer struct{
	sq_num int
	pool *DataPool
	sig_stopWrite chan bool
	wg sync.WaitGroup
	num *int
}

type Consumer struct{
	sq_num int
	pool *DataPool
	sig_stopRead chan bool
	wg sync.WaitGroup
}

type IF interface {
	start()
	stop()
}

func (pro *Producer) start(){
	var ch chan <-int
	ch = pro.pool.buff
	pro.pool.wg.Add(1)
	pro.wg.Add(1)
	fmt.Printf("production process %d started\n", pro.sq_num)


	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		select {
		case <- pro.sig_stopWrite:
			fmt.Printf("production process %d channel done()\n", pro.sq_num)
			pro.pool.wg.Done()
			pro.wg.Done()
			return
		default:
			pro.pool.cond.L.Lock()
			if pro.pool.length < pro.pool.size {
				*(pro.num)++
				ch <- *(pro.num)
				fmt.Printf("production process %d produce: %d\n", pro.sq_num, *(pro.num))
				pro.pool.length++
				
			} else {
				fmt.Printf("production process %d: DataPool is full\n",pro.sq_num)
			}
			pro.pool.cond.L.Unlock()
		}
	}
}

func (pro *Producer) stop(){
	close(pro.sig_stopWrite)
	pro.wg.Wait()
	fmt.Printf("production process %d exited\n", pro.sq_num)
}

func (con *Consumer) start(){
	var ch <-chan int
	ch = con.pool.buff
	con.pool.wg.Add(1)
	con.wg.Add(1)


	fmt.Printf("consumption process %d started\n", con.sq_num)
	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {
		select {
		case <- con.sig_stopRead:
			fmt.Printf("consumption process %d channel done()\n", con.sq_num)
			con.pool.wg.Done()
			con.wg.Done()
			return
		default:
			con.pool.cond.L.Lock()
			if con.pool.length > 0 {
				fmt.Printf("consumption process %d consume: %d\n", con.sq_num, <- ch)
				con.pool.length--
				
			} else {
				fmt.Printf("consumption process %d: DataPool is empty\n", con.sq_num)
			}
			con.pool.cond.L.Unlock()
			
		}			
	}
}

func (con *Consumer) stop(){
	close(con.sig_stopRead)
	con.wg.Wait()
	fmt.Printf("consumption process %d exited\n", con.sq_num)
}


func main(){

	const POOL_SIZE int = 10
	const NUM_PRODUCTION int = 2
	const NUM_CONSUMPTION int = 3

	var num int = 0

	pool := DataPool{
		buff: make(chan int, POOL_SIZE),
		size: POOL_SIZE,
		length: 0,
		cond: sync.NewCond(&sync.Mutex{}),
		wg: sync.WaitGroup{},
	}

	arr_producers := [NUM_PRODUCTION]*Producer{}
	for i:=0;i<NUM_PRODUCTION;i++ {
		pro := new(Producer)
		pro.sq_num=i
		pro.pool = &pool
		pro.sig_stopWrite = make(chan bool, 1)
		pro.wg = sync.WaitGroup{}
		pro.num = &num 
		arr_producers[i]=pro

		go pro.start()
	}

	arr_consumers := [NUM_CONSUMPTION]*Consumer{}
	for j:=0;j<NUM_CONSUMPTION;j++ {
		con := new(Consumer)
		con.sq_num=j
		con.pool = &pool
		con.sig_stopRead = make(chan bool, 1)
		con.wg = sync.WaitGroup{}
		arr_consumers[j]=con

		go con.start()
	}

	time.Sleep(20 * time.Second)

	for i:=0;i<NUM_PRODUCTION;i++ {
		go arr_producers[i].stop()
	}

	for i:=0;i<NUM_CONSUMPTION;i++ {
		go arr_consumers[i].stop()
	}

	pool.wg.Wait()
	close(pool.buff)
	fmt.Printf("main process exited\n")


}


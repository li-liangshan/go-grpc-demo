package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/grpc-demo/client"
	"github.com/grpc-demo/server"
	"github.com/grpc-demo/util"
)

func go_rpc() {
	go server.Start()
	client.Start()
}

func main() {
	// server.RunTcp()
	// client.RunTcp()

	// web.Run()
	// web.Http2ServerRun()
	// web.StatsRun()
	// web.RunStreamServer()

	util.RunId()
}

// 1.验证：runtime.Gosched()是会让出cpu的
func go_sched_test() {
	// 造场景，设置为单核那么就只能是并发，因为go1.5版本之后，默认是多核了
	runtime.GOMAXPROCS(1)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("go")
		}
	}()

	for i := 0; i < 2; i++ {
		fmt.Println("index:", i)
		runtime.Gosched()
		fmt.Println("hello")
	}
}

func go_ok_runtime() {
	runtime.GOMAXPROCS(2)
	go ok()
	time.Sleep(1 * time.Second)
}

func ok() {
	i := 1
	for {
		i++
		// fmt.Println(i) // 执行io了，会结束
		//time.Sleep(1) // 休眠需等待了，会结束
	}
}

var quit = make(chan int)

func loop(id string) {
	for i := 1; i <= 100; i++ {
		fmt.Printf("id-> %s %d\n", id, i)
		if i == 5 {
			time.Sleep(time.Second)
		}
	}
	quit <- 0
}

func go_runtime_loop() {
	runtime.GOMAXPROCS(1)
	go loop("id__001:")
	go loop("id__002:")

	for i := 1; i <= 2; i++ {
		<-quit
	}
}

// channel 频率控制
// 在对channel进行读写的时候，go还提供了非常人性化的操作，就是对读写的频率控制，通过time.Ticke实现
func channel_ticker() {
	requests := make(chan int, 5)
	for i := 1; i < 5; i++ {
		requests <- i
	}
	close(requests)
	limiter := time.Tick(time.Second * 1)
	for req := range requests {
		<-limiter
		fmt.Println("requets", req, time.Now()) //执行到这里，需要隔1秒才继续往下执行，time.Tick(timer)上面已定义
	}
}

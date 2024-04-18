package routines

import (
	"fmt"
	"time"
)

func dataFeedLB(data chan<- int, requests int) {
	// mimic requests
	for i := 0; i < requests; i++ {
		data <- i
	}
}

func spawnWorkersLB(n int, data <-chan int) {
	for i := 0; i < n; i++ {
		go workerLB(i, data)
	}
}

func workerLB(workerId int, data <-chan int) {
	for x := range data {
		if x == 999999 {
			fmt.Printf("-------------------------------- \n")
			fmt.Printf("Worker ---> %d <---- received the last request: %d\n", workerId, x)
			fmt.Printf("-------------------------------- \n")
			time.Sleep(time.Second)
		} else {
			fmt.Printf("Worker %d received %d\n", workerId, x)
			time.Sleep(time.Second)
		}
	}
}

func LoadBalancer(n, requests int) {
	data := make(chan int)

	spawnWorkersLB(n, data)
	dataFeedLB(data, requests)

	// close(data)
}

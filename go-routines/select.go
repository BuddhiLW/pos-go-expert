package routines

import "time"

func workerSelectPS(c chan int, v int) {
	// timeout := time.NewTimer(1 * time.Second)
	time.Sleep(1 * time.Second)
	// Channel `c` receives `v` value
	c <- v
}

func SelectPubSub() {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	go func() {
		c3 <- 3
		time.Sleep(1 * time.Second)
	}()
	go workerSelectPS(c1, 1)
	go workerSelectPS(c2, 2)

	for i := 0; i < 3; i++ {
		select {
		case msg1 := <-c1:
			println("received", msg1)
		case msg2 := <-c2:
			println("received", msg2)
		case msg3 := <-c3:
			println("received", msg3)
		case <-time.After(1 * time.Second):
			println("timeout")
			// default:
			// 	println("default")
		}
	}
}

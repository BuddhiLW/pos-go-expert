package routines

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var Global_count int = 0
var Global_count_A int = 0
var Global_count_B int = 0

func worker(id string, duration int) {
	fmt.Printf("Worker %s starting\n", id)

	// Sleep to simulate an expensive task.
	time.Sleep(time.Duration(duration) * time.Millisecond)
	fmt.Printf("Worker %s done\n\n", id)

	if id == "A" {
		Global_count_A++
	} else {
		Global_count_B++
	}
}

func task(name string, wg *sync.WaitGroup) {
	for i := 1; i < 11; i++ {
		fmt.Println("Global_count before task:", Global_count)
		fmt.Printf("%d: Task %s is running \n", i, name)
		if name == "A" {
			worker("A", 2)
		} else {
			worker("B", 1)
		}
		Global_count++
		fmt.Println("Global_count after task:", Global_count)
		fmt.Printf("Global count A: %d; Global count B: %d\n", Global_count_A, Global_count_B)

		wg.Done()
	}
}

func WaitGroups() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(10)

	go task("A", &waitGroup)
	go task("B", &waitGroup)

	// go func() {
	// 	for i := 1; i < 11; i++ {
	// 		fmt.Printf("%d: Task C is running \n", i)
	// 		time.Sleep(1 * time.Millisecond)
	// 		Global_count++
	// 		waitGroup.Done()
	// 	}
	// }()

	waitGroup.Wait()
}

var number uint64 = 0

func HTTP_race() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		number++
		log.Printf("Number: %d", number)
		w.Write([]byte(fmt.Sprintf("Number: %d", number)))
	})
	http.ListenAndServe(":3000", nil)
}

func HTTP_race_mux() {
	m := sync.Mutex{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		number++
		m.Unlock()

		log.Printf("Number: %d", number)
		w.Write([]byte(fmt.Sprintf("Number: %d", number)))
	})
	http.ListenAndServe(":3000", nil)
}

func Atomic() {
	// m := sync.Mutex{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Switch this code:
		// m.Lock()
		// number++
		// m.Unlock()

		// For this:
		atomic.AddUint64(&number, 1)
		log.Printf("Number: %d", number)
		w.Write([]byte(fmt.Sprintf("Number: %d", number)))
	})
	http.ListenAndServe(":3000", nil)
}

// Thread 1
func Channels(n int) {
	channel := make(chan string)

	var count int = 0
	// Thread 2
	go func() {
		for i := 1; i < n+1; i++ {
			var t time.Time
			if i == 1 {
				channel <- "Hello world! My first message in this Go Channel!"
				t = time.Now()
				fmt.Println("Thread 2 (pub), started filling channel, at the first time at:", t.Format("15:04:05.000000000"))
				count++

				time.Sleep(10 * time.Millisecond)
			} else {
				t = time.Now()
				fmt.Printf("Thread 2 (pub), filling channel again (%d time) at: %s\n", i, t.Format("15:04:05.000000000"))
				channel <- "Hello, filling out the channel again, after 2 seconds!"
				count++

				time.Sleep(10 * time.Millisecond)
			}
		}
		close(channel)
	}()

	// Pub -> competitive Sub
	for i := range channel {
		go func() {
			fmsg := <-channel
			fmt.Println(fmsg + " (sub from Thread 3)\n")
			count++
		}()

		// Thread 4
		go func() {
			fmsg := <-channel
			fmt.Println(fmsg + " (sub from Thread 4)\n")
			count++
		}()

		// Thread 1
		msg := i
		fmt.Println(msg + " (sub from Thread 1)\n")
	}

	// Pub -> []subs
	// for i := range channel {
	// 	go func() {
	// 		fmsg := i
	// 		fmt.Println(fmsg + " (sub from Thread 3)\n")
	// 		count++
	// 	}()

	// 	// Thread 4
	// 	go func() {
	// 		fmsg := i
	// 		fmt.Println(fmsg + " (sub from Thread 4)\n")
	// 		count++
	// 	}()

	// 	// // Thread 1
	// 	// msg := i
	// 	// fmt.Println(msg + " (sub from Thread 1)\n")
	// }

}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
	close(ch)
}

func subscription(ch chan int) {
	for i := range ch {
		fmt.Println(i)
	}
}

func PubSub() {
	ch := make(chan int)
	go publish(ch)
	subscription(ch)
}

// func readerPubSubWG(channel chan string, count int, wg *sync.WaitGroup) {
// 	// Pub -> competitive Sub
// 	for i := range channel {
// 		go func() {
// 			fmsg := <-channel
// 			fmt.Println(fmsg + " (sub from Thread 3)\n")
// 			count++
// 			wg.Done()
// 		}()

// 		// Thread 4
// 		go func() {
// 			fmsg := <-channel
// 			fmt.Println(fmsg + " (sub from Thread 4)\n")
// 			count++
// 			wg.Done()
// 		}()

// 		// Thread 1
// 		msg := i
// 		fmt.Println(msg + " (sub from Thread 1)\n")
// 		wg.Done()
// 	}
// }

// func publisherPubSubWG(channel chan string, n, count int) {
// 	for i := 1; i < n+1; i++ {
// 		var t time.Time

// 		if i == 1 {
// 			channel <- "Hello world! My first message in this Go Channel!"
// 			t = time.Now()
// 			fmt.Println("Thread 2 (pub), started filling channel, at the first time at:", t.Format("15:04:05.000000000"))
// 			count++

// 			time.Sleep(50 * time.Millisecond)
// 			// wg.Done()
// 		} else {
// 			t = time.Now()
// 			fmt.Printf("Thread 2 (pub), filling channel again (%d time) at: %s\n", i, t.Format("15:04:05.000000000"))
// 			channel <- "Hello, filling out the channel again, after 2 seconds!"
// 			count++

// 			time.Sleep(50 * time.Millisecond)
// 			// wg.Done()
// 		}
// 	}
// 	close(channel)
// }

// func WaitGroupPubSub(n int) {
// 	// Thread 1
// 	channel := make(chan string)
// 	wg := sync.WaitGroup{}
// 	wg.Add(n)

// 	var count int = 0
// 	go publisherPubSubWG(channel, n, count)

// 	// Thread 1
// 	readerPubSubWG(channel, count, &wg)
// 	wg.Wait()
// }

func subWG(id int, ch <-chan int, wg *sync.WaitGroup) {
	for x := range ch {
		fmt.Println("Thread " + strconv.Itoa(id) + " received: " + strconv.Itoa(x))
		wg.Done()
	}
}

func publishWG(ch chan<- int, n int) {
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
}

func WaitGroupPubSub(n int) {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(n)

	go publishWG(ch, n)
	go subWG(1, ch, &wg)

	wg.Wait()
}

// -----------------

package routines

import (
	"fmt"
	"log"
	"net/http"
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
		for i := 0; i < n; i++ {
			var t time.Time
			if i == 0 {
				channel <- "Hello world! My first message in this Go Channel!"
				t = time.Now()
				fmt.Println("Thread 2 (pub), started filling channel, at the first time at:", t.Format("15:04:05.000000000"))
				count++

				time.Sleep(2 * time.Second)
			} else {
				t = time.Now()
				fmt.Printf("Thread 2 (pub), filling channel again (%d time) at: %s\n", i, t.Format("15:04:05.000000000"))
				channel <- "Hello, filling the channel again, after 2 seconds!"
				count++

				time.Sleep(2 * time.Second)
			}
		}

	}()

	for count < n {
		// Thread 3
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
		msg := <-channel
		fmt.Println(msg + " (sub from Thread 1)\n")
	}
}

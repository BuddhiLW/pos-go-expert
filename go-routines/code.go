package routines

import (
	"fmt"
	"sync"
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

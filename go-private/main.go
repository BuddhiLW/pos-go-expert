package main

import (
	"fmt"

	"github.com/BuddhiLW/fcutils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}

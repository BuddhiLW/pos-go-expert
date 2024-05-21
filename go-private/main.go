package main

import (
	"fmt"

	events "github.com/BuddhiLW/fcutils-secret/pkg/events"
)

func main() {
	ed := events.NewEventDispatcher()
	fmt.Println(ed)
}

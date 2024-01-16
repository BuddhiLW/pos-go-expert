package main

import (
	"fmt"

	fundacao "github.com/buddhilw/pos-go-expert/fundacao"
)

func Spacing() {
	fmt.Println("\n--------------------\n")
}

func main() {
	fmt.Println("Structs:")
	fundacao.Structs()
	Spacing()

	fmt.Println("Interfaces:")
	fundacao.Interfaces()
	Spacing()

	fmt.Println("Pointers:")
	fundacao.Pointers()
	Spacing()

	fmt.Println("Modules:")
	fundacao.Modules()
}

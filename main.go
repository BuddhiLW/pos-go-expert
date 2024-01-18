package main

import (
	"fmt"

	fundacao "github.com/buddhilw/pos-go-expert/fundacao"
	packages "github.com/buddhilw/pos-go-expert/important-packages"
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
	Spacing()

	fmt.Println("File Manipulation:")
	packages.FileManipulation()
}

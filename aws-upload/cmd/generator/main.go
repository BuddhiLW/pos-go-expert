package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 0; i < 400; i++ {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.WriteString("Some random content")
	}
}

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	for i := 0; i < 1000; i++ {
		text := fmt.Sprintf("./tmp/file%d.txt", i)
		log.Print(text)
		f, err := os.Create(text)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		content := "random text in: " + text
		f.WriteString(content)
	}
}

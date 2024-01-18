package packages

import (
	"bufio"
	"fmt"
	"os"
)

func FileManipulation() {

	fmt.Println("\n# Create a File:")
	// create a file
	arq, err := os.Create("./sample-dir/new-file.txt")
	if err != nil {
		panic(err)
	}

	size, err := arq.Write([]byte("Hello World"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successefully created a file of size: %d bytes, in the directory %s \n", size, arq.Name())
	arq.Close()

	// read the file
	fmt.Println("\n# Read the Created File:")
	f, err := os.ReadFile("./sample-dir/new-file.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Reading the file...:", string(f))

	// reading a file by parts
	arq, err = os.Open("./sample-dir/new-file.txt")
	if err != nil {
		panic(err)
	}
	defer arq.Close()

	fmt.Println("\n# Read, using a streaming-buffer methodology:")
	stream := bufio.NewReader(arq)
	buffer := make([]byte, 5)
	fmt.Println("Start reading the file by parts:")
	for {
		// Allocate parts of the (file content) stream into the buffer
		n, err := stream.Read(buffer)

		// Throw an error if the end of the file is reached or if there is any other error
		// This will eventually break the loop
		if err != nil {
			fmt.Println(err)
			break
		}

		// Print the content of the buffer
		fmt.Println("Reading the file by parts...:", string(buffer[:n]))
	}

	fmt.Println("\n# Deleting the Created File:")
	err = os.Remove("./sample-dir/new-file.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Successefully deleted the file.")
}

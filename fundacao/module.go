package fundacao

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

func ReadDirRecursive(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {

		if e.IsDir() {
			newDir := dir + e.Name() + "/"
			// fmt.Println("Inside new dir:", newDir)
			ReadDirRecursive(newDir)
		} else {
			// fmt.Println("Inside dir:", dir)
			filep := dir + e.Name()
			fmt.Println("File: ", filep)
		}
	}
}

func Modules() {
	fmt.Println("(using uuid library) new uid:", uuid.New())

	baseDir := "./sample-dir/"
	fmt.Println("Using =os= library to read (non-empty) directories recusively:")
	fmt.Println("Reading from the root directory:", baseDir)
	ReadDirRecursive(baseDir)
}

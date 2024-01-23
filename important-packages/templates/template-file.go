package templates

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

type Courses []Course

func Wd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}

func ExternalFileTemplate() {
	fmt.Println("Templating using external file")
	fileName := "courses.html"
	file := Wd() + "/templates/" + fileName
	fmt.Println("Template file directory: " + file)
	r, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Template file contents: " + string(r))

	t := template.Must(template.New(fileName).ParseFiles(file))
	err = t.Execute(os.Stdout, Courses{
		{"Go Expert", 360},
		{"Go Expert II", 720},
		{"Go Expert III", 1080},
	})
	if err != nil {
		panic(err)
	}
}

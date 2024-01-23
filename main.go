package main

import (
	"fmt"
	"runtime"

	"github.com/buddhilw/pos-go-expert/fundacao"
	http_funcs "github.com/buddhilw/pos-go-expert/http-funcs"
	packages "github.com/buddhilw/pos-go-expert/important-packages"
	templates "github.com/buddhilw/pos-go-expert/important-packages/templates"
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
	Spacing()

	fmt.Println("HTTP:")
	packages.CEPSearch("")
	Spacing()

	fmt.Println("JSON:")
	packages.Json()
	Spacing()

	fmt.Println("HTTP Server for searching CEPs: (localhost:8989)")
	go http_funcs.CEP()
	Spacing()

	fmt.Println("HTTP Server for serving files: (localhost:8990)")
	go http_funcs.FileServer()
	Spacing()

	defer runtime.Goexit()

	fmt.Println("Templating system:")
	templates.First()
	templates.Must()
	templates.ExternalFileTemplate()
	Spacing()

	fmt.Println("Template Web Server: (localhost:8080)")
	go templates.TemplateWebServer()
	Spacing()

	// defer fmt.Println("Server exiting!")
}

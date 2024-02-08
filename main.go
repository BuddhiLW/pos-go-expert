package main

import (
	"fmt"
)

func Spacing() {
	fmt.Println("\n--------------------\n")
}

func main() {
	// fmt.Println("Structs:")
	// fundacao.Structs()
	// Spacing()

	// fmt.Println("Interfaces:")
	// fundacao.Interfaces()
	// Spacing()

	// fmt.Println("Pointers:")
	// fundacao.Pointers()
	// Spacing()

	// fmt.Println("Modules:")
	// fundacao.Modules()
	// Spacing()

	// fmt.Println("File Manipulation:")
	// packages.FileManipulation()
	// Spacing()

	// fmt.Println("HTTP:")
	// packages.CEPSearch("")
	// Spacing()

	// fmt.Println("JSON:")
	// packages.Json()
	// Spacing()

	// fmt.Println("HTTP Server for searching CEPs: (localhost:8989)")
	// go http_funcs.CEP()
	// Spacing()

	// fmt.Println("HTTP Server for serving files: (localhost:8990)")
	// go http_funcs.FileServer()
	// Spacing()

	// defer runtime.Goexit()

	// fmt.Println("Templating system:")
	// templates.First()
	// templates.Must()
	// templates.ExternalFileTemplate()
	// Spacing()

	// fmt.Println("Template Web Server: (localhost:8080)")
	// go templates.TemplateWebServer()
	// Spacing()

	// fmt.Println("Calling Google with Timeout")
	// go http_funcs.TimeOutHTTP()
	// Spacing()

	// fmt.Println("Posting to Google")
	// go http_funcs.PostHTTP()
	// Spacing()

	// fmt.Println("Context - Will it book? - example:")
	// context.Context()
	// Spacing()

	fmt.Println("DB - New product:")
	db.NewProductExample()
	Spacing()

	fmt.Println("DB - Search by product-name:")
	db.SearchByNameExample()
	Spacing()

	fmt.Println("DB - Update product-price by product-name:")
	db.UpdateProductExample()
	db.SearchByNameExample()
	Spacing()

	fmt.Println("DB - Update product-price by product-id:")
	db.UpdateProductExample2()
	Spacing()

	fmt.Println("DB - All products:")
	db.AllProductsExample()
	Spacing()

	fmt.Println("DB - Delete all products:")
	db.DeleteAllProductsExample()
	db.AllProductsExample()
	Spacing()
}

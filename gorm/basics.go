package exp_gorm

import (
	"fmt"

	"gorm.io/gorm"
)

func GormCreateProductsExample() {
	db := ConnectDB()
	Migration(db)
	// Create
	products := []Product{
		{Name: "Laptop", Price: 1000},
		{Name: "Mouse", Price: 100},
		{Name: "Keyboard", Price: 250},
		{Name: "Monitor", Price: 500},
	}
	db.Create(&products)
}

func GormReadProductsExample() {
	db := ConnectDB()

	// Read
	// fmt.Println("\nSearch first product (with id >= 1):")
	var product Product
	db.Where("id >= 1").First(&product) // find product with id 1
	fmt.Println(product)

	fmt.Println("\nSearch product with name 'Keyboard':")
	var product2 Product
	db.First(&product2, "name = ?", "Keyboard") // find product with name Laptop
	fmt.Println(product2)

	fmt.Println("\nSearch all products:")
	var products []Product
	db.Find(&products)
	fmt.Println(products)
}

func GormDeleteAllProductsExample() {
	db := ConnectDB()
	// https://gorm.io/docs/delete.html#Block-Global-Delete
	// If you perform a batch delete without any conditions, GORM WON’T run it, and will return ErrMissingWhereClause error
	// You have to use some conditions or use raw SQL or enable AllowGlobalUpdate mode, for example:
	// db.Where("1 = 1").Delete(&Product{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Product{})
	// db.Delete(&Product{}).E
}

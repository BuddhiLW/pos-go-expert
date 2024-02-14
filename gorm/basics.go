package exp_gorm

import (
	"fmt"

	"gorm.io/gorm"
)

func GormCreateProductsExample() {
	db := ConnectDB()
	Migration(db)

	var categories []Category
	db.Find(&categories)
	fmt.Println(categories)

	// db.Create(&SerialNumber{Number: "1234", ProductID: 1})

	// Create
	products := []Product{
		{Name: "Laptop", Price: 1000,
			Categories: []Category{categories[0], categories[3]}},
		{Name: "Mouse", Price: 100,
			Categories: []Category{categories[0], categories[1]}},
		{Name: "Keyboard", Price: 250,
			Categories: []Category{categories[0], categories[2]}},
		{Name: "Monitor", Price: 500,
			Categories: []Category{categories[0], categories[1], categories[2]}}}
	db.Create(&products)
}

func PaginationExample() {
	db := ConnectDB()
	// Pagination template
	var products []Product
	db.Limit(2).Offset(2).Find(&products)
	fmt.Println(products)
}

func RegexSearchExample() {
	db := ConnectDB()
	var products []Product
	db.Where("name LIKE ?", "%"+"board"+"%").Find(&products)
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

func GormUpdateProductsExample() {
	db := ConnectDB()
	// Update - update product's price to 200
	db.Model(&Product{}).Where("name = ?", "Mouse").Update("Price", 200)

	// Alternative
	var p Product
	db.Where("name = ?", "Mouse").First(&p)
	p.Name = "Mouse Updated"
	db.Save(&p)
}

func GormDeleteProductExample() {
	db := ConnectDB()
	// Delete - delete product with lowest id
	var product Product
	db.Where("id >= 1").First(&product) // find product with id 1
	db.Delete(&product)
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

	fmt.Println("\nPrototype of pagination:")
	PaginationExample()

	fmt.Println("\nPrototype of regex search:")
	RegexSearchExample()

}

func GormUpdateExamples() {
	fmt.Println("\nUpdate examples -- Mouse:")
	GormUpdateProductsExample()

	db := ConnectDB()
	fmt.Println("\nSearch all products:")
	var products []Product
	db.Find(&products)
	fmt.Println(products)
}

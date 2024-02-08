package db

import (
	"database/sql"
	"fmt"
)

func UpdateProductByName(name string, price float64) error {
	db := ConnectDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE products SET price = ? WHERE name = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(price, name)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProductById(db *sql.DB, product *Product) error {

	stmt, err := db.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProductExample() {
	err := UpdateProductByName("new product", 20.0)
	if err != nil {
		panic(err)
	}
}

func UpdateProductExample2() {
	newProduct := NewProduct("Id Product", 10.0)
	fmt.Println("The new product:", newProduct)
	con := ConnectDB()
	InsertProductDB(con, newProduct)

	newProduct.Name = "Updated Product Name, by ID"
	newProduct.Price = 100.0

	err := UpdateProductById(con, newProduct)
	if err != nil {
		panic(err)
	}
}

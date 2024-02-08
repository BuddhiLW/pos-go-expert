package db

import (
	"database/sql"

	"github.com/google/uuid"
)

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func InsertProductDB(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("INSERT INTO products(id, name, price) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func NewProductExample() {
	db := ConnectDB()
	defer db.Close()

	product := NewProduct("new product", 10.0)
	err := InsertProductDB(db, product)
	if err != nil {
		panic(err)
	}
}

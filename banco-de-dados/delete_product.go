package db

import (
	"database/sql"
	"fmt"
)

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func deleteAllProducts(db *sql.DB) error {
	allProducts, err := AllProducts()
	for _, product := range allProducts {
		err = deleteProduct(db, product.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteAllProductsExample() {
	db := ConnectDB()
	defer db.Close()

	err := deleteAllProducts(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("All products deleted, successfully!")
}

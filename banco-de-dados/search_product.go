package db

func SearchByName(name string) (*Product, error) {
	db := ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM  products WHERE name = ?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var product Product

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			panic(err)
		}
	}

	return &product, nil
}

func PrintProduct(product *Product) {
	println("ID:", product.ID)
	println("Name:", product.Name)
	println("Price:", product.Price)
}

func SearchByNameExample() {
	product, err := SearchByName("new product")
	if err != nil {
		panic(err)
	}
	PrintProduct(product)
}

func AllProducts() ([]Product, error) {
	db := ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func AllProductsExample() {
	products, err := AllProducts()
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		PrintProduct(&product)
	}
}

func SearchById(id string) (*Product, error) {
	db := ConnectDB()
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, price from products where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var product Product
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

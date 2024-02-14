package exp_gorm

// func PreloadedFind() []Product {
// 	var products []Product
// 	db := ConnectDB()
// 	db.Preload("Category").Preload("SerialNumber").Find(&products)
// 	return products
// }

// func PreloadedFindExample() {
// 	products := PreloadedFind()
// 	for _, product := range products {
// 		fmt.Println(product.Name, product.Category.Name, product.SerialNumber.Number)
// 	}
// 	product := products[0]
// 	fmt.Println(product.Name, product.Category.Name, product.SerialNumber.Number)
// }

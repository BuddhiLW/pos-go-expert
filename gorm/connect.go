package exp_gorm

import (
	// _ "gorm.io/driver/mysql"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID         int `gorm:"primary_key"`
	Name       string
	Price      float64
	CategoryID int `gorm:"foreignkey:ID;default:1"`
	Category   *Category
	gorm.Model
}

type Category struct {
	ID   int `gorm:"primary_key"`
	Name string
}

func ConnectDB() *gorm.DB {
	// -----
	//
	// sqlDB, err := sql.Open("mysql", "mydb_dsn")
	// gormDB, err := gorm.Open(mysql.New(mysql.Config{
	//   Conn: sqlDB,
	// }), &gorm.Config{})

	dns := "buddhilw:pass@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func Migration(db *gorm.DB) {
	db.AutoMigrate(&Product{}, &Category{})

	var c0 Category
	db.First(&c0, 1)
	if c0.ID == 1 {
		return
	} else {
		category := Category{Name: "Generic"}
		db.Create(&category)
	}
}

func AutoMigrateExample() {
	db := ConnectDB()
	Migration(db)
}

// func main() {
// 	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()

// 	// Migrate the schema
// 	db.AutoMigrate(&Product{})

// 	// Create
// 	db.Create(&Product{Name: "Laptop", Price: 1000})

// 	// Read
// 	var product Product
// 	db.First(&product, 1)                    // find product with id 1
// 	db.First(&product, "name = ?", "Laptop") // find product with name Laptop

// 	// Update - update product's price to 2000
// 	db.Model(&product).Update("Price", 2000)

// 	// Delete - delete product
// 	db.Delete(&product)
// }

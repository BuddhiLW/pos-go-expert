package exp_gorm

import "gorm.io/gorm"

type Product struct {
	ID    int `gorm:"primary_key"`
	Name  string
	Price float64
	// CategoryID int `gorm:"foreignkey:ID;default:1"`
	// Category     *Category
	Categories []Category `gorm:"many2many:product_categories;"`
	// SerialNumber *SerialNumber
	gorm.Model
}

type Category struct {
	ID       int `gorm:"primary_key"`
	Name     string
	Products []Product `gorm:"many2many:product_categories;"`
}

// type SerialNumber struct {
// 	ID        int `gorm:"primary_key"`
// 	Number    string
// 	ProductID int
// }

// type CategoryM2M struct {
// 	ID   int `gorm:"primary_key"`
// 	Name string
// }

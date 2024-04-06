package database

import "github.com/BuddhiLW/pos-go-expert/APIs/internal/entity"

type UserInterface interface {
	Createuser(user *entity.User)
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}

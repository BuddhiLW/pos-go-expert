package database

import "github.com/BuddhiLW/pos-go-expert/APIs/internal/entity"

type UserInterface interface {
	Createuser(user *entity.User)
	FindByEmail(email string) (*entity.User, error)
}

package entity

import (
	"errors"
	"time"

	"github.com/BuddhiLW/pos-go-expert/APIs/pkg/entity"
)

var (
	ErrIdRequired    = errors.New("ID is required")
	ErrInvalidID     = errors.New("ID is invalid")
	ErrNameRequired  = errors.New("Name is required")
	ErrPriceRequired = errors.New("Prince is required")
	ErrInvalidPrice  = errors.New("Price is invalid")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:        entity.NewID(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}

	err := product.Validate()
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIdRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price == 0 {
		return ErrPriceRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

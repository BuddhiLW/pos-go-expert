package dto

type CreateProductRequest struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

// type UpdateProductRequest struct {
// 	Name  string  `json:"name"`
// 	Price float64 `json:"price"`
// }

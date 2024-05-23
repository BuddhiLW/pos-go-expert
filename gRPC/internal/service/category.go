package service

import (
	"context"

	"github.com/BuddhiLW/pos-go-expert/gRPC/internal/database"
	"github.com/BuddhiLW/pos-go-expert/gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{CategoryDB: categoryDB}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method CreateCategory not implemented")
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{Category: categoryResponse}, nil
}

// CreateCategory(context.Context, *CreateCategoryRequest) (*CategoryResponse, error)

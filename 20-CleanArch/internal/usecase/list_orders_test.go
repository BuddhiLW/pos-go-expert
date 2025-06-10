package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ordersystem/internal/entity"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *entity.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetAll() ([]*entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func TestListOrdersUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		mockOrders     []*entity.Order
		mockError      error
		expectedOrders int
		expectedError  bool
	}{
		{
			name: "successful retrieval with multiple orders",
			mockOrders: []*entity.Order{
				{ID: "1", Price: 100, Tax: 10, FinalPrice: 110},
				{ID: "2", Price: 200, Tax: 20, FinalPrice: 220},
				{ID: "3", Price: 50, Tax: 5, FinalPrice: 55},
			},
			mockError:      nil,
			expectedOrders: 3,
			expectedError:  false,
		},
		{
			name:           "successful retrieval with empty list",
			mockOrders:     []*entity.Order{},
			mockError:      nil,
			expectedOrders: 0,
			expectedError:  false,
		},
		{
			name: "successful retrieval with single order",
			mockOrders: []*entity.Order{
				{ID: "1", Price: 100, Tax: 10, FinalPrice: 110},
			},
			mockError:      nil,
			expectedOrders: 1,
			expectedError:  false,
		},
		{
			name:           "repository error",
			mockOrders:     nil,
			mockError:      errors.New("database connection error"),
			expectedOrders: 0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockOrderRepository)
			mockRepo.On("GetAll").Return(tt.mockOrders, tt.mockError)

			useCase := NewListOrdersUseCase(mockRepo)

			// Act
			result, err := useCase.Execute()

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, ListOrdersOutputDTO{}, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result.Orders, tt.expectedOrders)

				// Verify DTO conversion
				for i, order := range tt.mockOrders {
					assert.Equal(t, order.ID, result.Orders[i].ID)
					assert.Equal(t, order.Price, result.Orders[i].Price)
					assert.Equal(t, order.Tax, result.Orders[i].Tax)
					assert.Equal(t, order.FinalPrice, result.Orders[i].FinalPrice)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestListOrdersUseCase_Execute_DTOConversion(t *testing.T) {
	// Arrange
	mockOrders := []*entity.Order{
		{ID: "order-1", Price: 99.99, Tax: 9.99, FinalPrice: 109.98},
		{ID: "order-2", Price: 199.50, Tax: 19.95, FinalPrice: 219.45},
	}

	mockRepo := new(MockOrderRepository)
	mockRepo.On("GetAll").Return(mockOrders, nil)

	useCase := NewListOrdersUseCase(mockRepo)

	// Act
	result, err := useCase.Execute()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result.Orders, 2)

	// Verify first order DTO
	assert.Equal(t, "order-1", result.Orders[0].ID)
	assert.Equal(t, 99.99, result.Orders[0].Price)
	assert.Equal(t, 9.99, result.Orders[0].Tax)
	assert.Equal(t, 109.98, result.Orders[0].FinalPrice)

	// Verify second order DTO
	assert.Equal(t, "order-2", result.Orders[1].ID)
	assert.Equal(t, 199.50, result.Orders[1].Price)
	assert.Equal(t, 19.95, result.Orders[1].Tax)
	assert.Equal(t, 219.45, result.Orders[1].FinalPrice)

	mockRepo.AssertExpectations(t)
}

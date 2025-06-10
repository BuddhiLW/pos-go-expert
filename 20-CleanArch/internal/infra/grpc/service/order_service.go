package service

import (
	"context"

	"ordersystem/internal/entity"
	"ordersystem/internal/infra/grpc/pb"
	"ordersystem/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	OrderRepository    entity.OrderRepositoryInterface
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, orderRepository entity.OrderRepositoryInterface) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		OrderRepository:    orderRepository,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	listOrders := usecase.NewListOrdersUseCase(s.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		return nil, err
	}

	var grpcOrders []*pb.CreateOrderResponse
	for _, order := range output.Orders {
		grpcOrder := &pb.CreateOrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}
		grpcOrders = append(grpcOrders, grpcOrder)
	}

	return &pb.ListOrdersResponse{
		Orders: grpcOrders,
	}, nil
}

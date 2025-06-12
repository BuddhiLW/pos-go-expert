package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"ordersystem/configs"
	"ordersystem/internal/event/handler"
	"ordersystem/internal/infra/database"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"ordersystem/internal/infra/graph"
	"ordersystem/internal/infra/grpc/pb"
	"ordersystem/internal/infra/grpc/service"
	"ordersystem/internal/infra/web/webserver"
	"ordersystem/pkg/events"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	// Use Wire-generated dependency injection
	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	// Create repository for gRPC service
	orderRepository := database.NewOrderRepository(db)
	createOrderService := service.NewOrderService(*createOrderUseCase, orderRepository)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
		OrderRepository:    orderRepository,
	}}))

	// Create a separate mux for GraphQL
	graphqlMux := http.NewServeMux()
	graphqlMux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	graphqlMux.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	go http.ListenAndServe(":"+configs.GraphQLServerPort, graphqlMux)

	// Keep the main goroutine alive
	select {}
}

func getRabbitMQChannel(configs *configs.Conf) *amqp.Channel {
	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		configs.RabbitMQUser,
		configs.RabbitMQPassword,
		configs.RabbitMQHost,
		configs.RabbitMQPort)

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

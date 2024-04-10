package main

import (
	"fmt"
	"net/http"

	"github.com/BuddhiLW/pos-go-expert/APIs/configs"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/entity"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/infra/database"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := configs.LoadConf(".")
	// fmt.Println(config)

	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/products", productHandler.CreateProduct)
	// http.HandleFunc("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)

	fmt.Println("Serving at:", string(config.WebSeverPort))
	http.ListenAndServe(":"+string(config.WebSeverPort), r)
}

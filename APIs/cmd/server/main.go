package main

import (
	"fmt"
	"net/http"

	"github.com/BuddhiLW/pos-go-expert/APIs/configs"
	_ "github.com/BuddhiLW/pos-go-expert/APIs/docs"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/entity"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/infra/database"
	"github.com/BuddhiLW/pos-go-expert/APIs/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Pedro Gomes Branquinho
// @contact.url    http://www.buddhilw.com
// @contact.email  pedrogbranquinho@gmail.com

// @license.name   Full Cycle License
// @license.url    http://www.fullcycle.com.br

// @host      localhost:8009
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	config, err := configs.LoadConf(".")

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
	// r.Use(LogRequest)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8009/docs/doc.json")))

	fmt.Println("Serving at:", string(config.WebSeverPort))
	http.ListenAndServe(":"+string(config.WebSeverPort), r)
}

// Como um middleware looks like
// func LogRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			log.Printf("Request: %s %s", r.Method, r.URL.Path)
// 			next.ServeHTTP(w, r)
// 		})
// }

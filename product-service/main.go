package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	deliveryHttp "github.com/sinulingga23/microservices/product-service/delivery/http"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/service"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	client, err := repository.ConnectMongo()
	if err != nil {
		log.Printf("Failed connect to mongo: %v", err)
		return
	}

	categoryRepository := repository.NewCategoryRepository(client)
	productRepository := repository.NewProductRepository(client)

	categoryService := service.NewCategoryService(categoryRepository)
	productService := service.NewProductService(categoryRepository, productRepository)

	categoryHandler := deliveryHttp.NewCategoryHandler(categoryService)
	productHandler := deliveryHttp.NewProductHandler(productService)

	categoryHandler.BindRoutes(r)
	productHandler.BindRoutes(r)

	log.Println("Running product-service on :8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Printf("Error listen and serve: %v", err)
		return
	}
}

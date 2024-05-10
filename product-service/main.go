package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	deliveryHttp "github.com/sinulingga23/microservices/product-service/delivery/http"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/service"
)

func main() {
	r := chi.NewRouter()

	client, err := repository.ConnectMongo()
	if err != nil {
		log.Printf("Failed connect to mongo: %v", err)
		return
	}

	categoryRepository := repository.NewCategoryRepository(client)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := deliveryHttp.NewCategoryHandler(categoryService)

	categoryHandler.BindRoutes(r)

	log.Println("Running product-service on :8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Printf("Error listen and serve: %v", err)
		return
	}
}

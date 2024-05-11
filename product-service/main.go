package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	deliveryGrpc "github.com/sinulingga23/microservices/product-service/delivery/grpc"
	deliveryHttp "github.com/sinulingga23/microservices/product-service/delivery/http"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/service"
	"github.com/sinulingga23/microservices/product-service/service/rpc"
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

	productRpc := rpc.NewProductRpc(productRepository)

	categoryHandler := deliveryHttp.NewCategoryHandler(categoryService)
	productHandler := deliveryHttp.NewProductHandler(productService)

	categoryHandler.BindRoutes(r)
	productHandler.BindRoutes(r)

	listener, err := net.Listen("tcp", ":3031") // for grpc
	if err != nil {
		log.Printf("Err listener: %v", err)
		return
	}
	grpcServer := deliveryGrpc.InitRegistrationServer(productRpc)

	// temporary solution
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("Running product-service on :8081")
		err = http.ListenAndServe(":8081", r)
		if err != nil {
			log.Printf("Error listen and serve: %v", err)
			return
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		log.Println("Running grpc-server on :3031")
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Printf("Failed to served grpc-server: %v", err)
			return
		}
	}(&wg)

	wg.Wait()
}

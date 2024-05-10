package service

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sinulingga23/microservices/product-service/model"
	"github.com/sinulingga23/microservices/product-service/payload"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type IProductService interface {
	HandleAddProduct(ctx context.Context, request payload.AddProductRequest) payload.ResponseGeneral
}

type productService struct {
	categoryRepository repository.ICategoryRepository
	productRepository  repository.IProductRepository
}

func NewProductService(
	categoryRepository repository.ICategoryRepository,
	productRepository repository.IProductRepository,
) IProductService {
	return &productService{categoryRepository: categoryRepository, productRepository: productRepository}
}

func (s *productService) HandleAddProduct(ctx context.Context, request payload.AddProductRequest) payload.ResponseGeneral {
	serviceName := "handle_add_product"
	response := payload.NewResponseGeneral(http.StatusOK, "Success")

	if err := request.Validate(); err != nil {
		log.Printf("%s: error validate: %v", serviceName, err)
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		return response
	}

	uniqueCategoryIds := []string{}
	visitedId := map[string]int{}
	for i := 0; i < len(request.CategoryIds); i++ {
		id := request.CategoryIds[i]

		_, ok := visitedId[id]
		if !ok {
			uniqueCategoryIds = append(uniqueCategoryIds, id)
		}

		visitedId[id] = 1
	}

	resultIds := []string{}
	var err error

	if len(uniqueCategoryIds) > 0 {
		resultIds, err = s.categoryRepository.GetIdsByIds(ctx, uniqueCategoryIds)
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		log.Printf("%s: category not found: %v", serviceName, err)
		response.StatusCode = http.StatusBadRequest
		response.Message = utils.ErrCategoryNotExists.Error()
		return response
	}
	if err != nil {
		log.Printf("%s: error when try get ids: %v", serviceName, err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = utils.ErrDBError.Error()
		return response
	}
	if len(resultIds) == 0 {
		response.StatusCode = http.StatusBadRequest
		response.Message = utils.ErrCategoryNotExists.Error()
		return response
	}

	count := 0
	for i := 0; i < len(resultIds); i++ {
		id := resultIds[i]
		visitedId[id] += 1
		if visitedId[id] == 2 {
			count += 1
		}
	}

	if count != len(resultIds) {
		response.StatusCode = http.StatusBadRequest
		response.Message = utils.ErrCategoryNotExists.Error()
		return response
	}

	err = s.productRepository.Create(ctx, model.Product{
		Id:          uuid.NewString(),
		Name:        request.Name,
		Qtty:        request.Qtty,
		Price:       request.Price,
		Description: request.Description,
		CategoryIds: resultIds,
	})
	if err != nil {
		log.Printf("%s: error when try create data: %v", serviceName, err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = utils.ErrDBError.Error()
		return response
	}

	return response
}

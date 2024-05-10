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
	if len(request.CategoryIds) > 0 {

		mapId := map[string]int{}
		for i := 0; i < len(request.CategoryIds); i++ {
			_, ok := mapId[request.CategoryIds[i]]
			if !ok {
				uniqueCategoryIds = append(uniqueCategoryIds, request.CategoryIds[i])
			}

			mapId[request.CategoryIds[i]] = 1
		}

		ids, err := s.categoryRepository.GetIdsByIds(ctx, uniqueCategoryIds)
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

		for i := 0; i < len(ids); i++ {
			mapId[ids[i]] += 1
		}

		isNotExists := false
		for _, ct := range mapId {
			if ct == 1 {
				isNotExists = true
				break
			}
		}

		if isNotExists {
			response.StatusCode = http.StatusBadRequest
			response.Message = utils.ErrCategoryNotExists.Error()
			return response
		}
	}

	err := s.productRepository.Create(ctx, model.Product{
		Id:          uuid.NewString(),
		Name:        request.Name,
		Qtty:        request.Qtty,
		Price:       request.Price,
		Description: request.Description,
		CategoryIds: uniqueCategoryIds,
	})
	if err != nil {
		log.Printf("%s: error when try create data: %v", serviceName, err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = utils.ErrDBError.Error()
		return response
	}

	return response
}

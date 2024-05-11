package service

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sinulingga23/microservices/product-service/model"
	"github.com/sinulingga23/microservices/product-service/payload"
	"github.com/sinulingga23/microservices/product-service/repository"
	"github.com/sinulingga23/microservices/product-service/utils"
)

type ICategoryService interface {
	HandleAddCategory(ctx context.Context, request payload.AddCategoryRequest) payload.ResponseGeneral
}

type categoryService struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryService(categoryRepository repository.ICategoryRepository) ICategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (s *categoryService) HandleAddCategory(ctx context.Context, request payload.AddCategoryRequest) payload.ResponseGeneral {
	serviceName := "handle_add_category"
	response := payload.NewResponseGeneral(http.StatusOK, "Success")

	if err := request.Validate(); err != nil {
		log.Printf("%s: Error Validate: %v", serviceName, err)
		response.StatusCode = http.StatusBadRequest
		response.Message = err.Error()
		return response
	}

	err := s.categoryRepository.Create(ctx, model.Category{
		Id:   uuid.NewString(),
		Name: request.Name,
	})
	if err != nil {
		log.Printf("%s: Error Create Data: %v", serviceName, err)
		response.StatusCode = http.StatusInternalServerError
		response.Message = utils.ErrDBError.Error()
		return response
	}

	return response
}

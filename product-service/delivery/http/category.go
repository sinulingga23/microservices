package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sinulingga23/microservices/product-service/payload"
	"github.com/sinulingga23/microservices/product-service/service"
)

type categoryHandler struct {
	categoryService service.ICategoryService
}

func NewCategoryHandler(categoryService service.ICategoryService) *categoryHandler {
	return &categoryHandler{categoryService: categoryService}
}

func (h *categoryHandler) BindRoutes(r *chi.Mux) {
	r.Post("/api/v1/categories", h.HandleAddCategory)
}

func (h *categoryHandler) HandleAddCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := payload.AddCategoryRequest{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := h.categoryService.HandleAddCategory(context.TODO(), request)
	bytesResponse, err := json.Marshal(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesResponse)
}

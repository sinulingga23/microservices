package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sinulingga23/microservices/product-service/payload"
	"github.com/sinulingga23/microservices/product-service/service"
)

type productHandler struct {
	productService service.IProductService
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{productService: productService}
}

func (h *productHandler) BindRoutes(r *chi.Mux) {
	r.Post("/api/v1/products", h.HandleAddProduct)
}

func (h *productHandler) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := payload.AddProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := h.productService.HandleAddProduct(context.TODO(), request)
	bytesResponse, err := json.Marshal(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write(bytesResponse)
}

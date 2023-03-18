package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	ProductRequest struct {
		id    string
		Name  string `json:"name"`
		Qtty  int    `json:"qtty"`
		Price int    `json:"price"`
	}

	Product struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		Qtty      int       `json:"qtty"`
		Price     int       `json:"price"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	ProductRepository struct {
		Items []Product
	}
)

var (
	port              = "8081"
	productRepository ProductRepository
	mu                sync.Mutex
	ErrRecordNotFound = errors.New("record not found")
)

func (p *ProductRepository) Add(productRequest ProductRequest) {

	p.Items = append(p.Items, Product{
		Id:        productRequest.id,
		Name:      productRequest.Name,
		Qtty:      productRequest.Qtty,
		Price:     productRequest.Price,
		CreatedAt: time.Now(),
	})
}

func (p *ProductRepository) FindProductById(id string) (Product, error) {

	items := p.Items
	lenItems := len(items)
	if lenItems == 0 {
		return Product{}, ErrRecordNotFound
	}

	for i := 0; i < lenItems; i++ {
		if id == items[i].Id {
			return items[i], nil
		}
	}

	return Product{}, ErrRecordNotFound
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productRequest := ProductRequest{}

	bytesBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		log.Printf("errReadAll: %v", errReadAll)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errUnmarshal := json.Unmarshal(bytesBody, &productRequest); errUnmarshal != nil {
		log.Printf("errUnmarshal: %v", errUnmarshal)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	productId, errGenerateProductId := generateProductId(len(productRepository.Items))
	if errGenerateProductId != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productRequest.id = productId
	productRepository.Add(productRequest)
	w.WriteHeader(http.StatusOK)
	return
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(productRepository.Items) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	bytesItems, errMarshal := json.Marshal(&productRepository.Items)
	if errMarshal != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(bytesItems))
	w.WriteHeader(http.StatusOK)
	return
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	product, errFindProductById := productRepository.FindProductById(id)
	if errFindProductById != nil {
		if errors.Is(errFindProductById, ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytesProduct, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(bytesProduct)
	w.WriteHeader(http.StatusOK)
	return
}

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Post("/api/v1/products", addProduct)
	router.Get("/api/v1/products", getProducts)
	router.Get("/api/v1/products/{id}", getProductById)

	log.Printf("Running product-service on :%s", port)
	if errListenAndServe := http.ListenAndServe(fmt.Sprintf(":%s", port), router); errListenAndServe != nil {
		log.Fatalf("Error when listen and serve: %v", errListenAndServe)
	}
}

func generateProductId(currentSize int) (string, error) {
	if currentSize <= -1 {
		return "", errors.New("currentSize should greater than -1")
	}

	currentSize += 1

	tempCurrentSize := currentSize
	countDigit := 0
	for currentSize != 0 {
		currentSize /= 10
		countDigit += 1
	}

	if tempCurrentSize == 0 {
		tempCurrentSize = 1
		countDigit = 1
	}

	totalDigitZero := 6
	totalDigitZero -= countDigit
	return fmt.Sprintf("PRD%s%d", strings.Repeat("0", totalDigitZero), tempCurrentSize), nil
}

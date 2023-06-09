package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinulingga23/microservices/product-service/constant"
	"github.com/sinulingga23/microservices/product-service/monitoring"
	"github.com/sinulingga23/microservices/product-service/utils"
)

type (
	ProductRequest struct {
		id    string
		Name  string `json:"name"`
		Stock int    `json:"stock"`
		Price int    `json:"price"`
	}

	Product struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		Stock     int       `json:"stock"`
		Price     int       `json:"price"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	Deduction struct {
		Id            string
		ProductId     string
		OrderId       string
		OrderDetailId string
		Qtty          int
	}

	ProductRepository struct {
		Items            []Product
		HistoryDeduction []Deduction
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
		Stock:     productRequest.Stock,
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

func (p *ProductRepository) FindProductsByIds(ids []string) ([]Product, error) {
	items := p.Items
	lenItems := len(items)
	if lenItems == 0 {
		return []Product{}, ErrRecordNotFound
	}

	products := make([]Product, 0)
	lenIds := len(ids)
	for i := 0; i < lenIds; i++ {
		for j := 0; j < lenItems; j++ {
			product := items[j]
			if ids[i] == product.Id {
				products = append(products, product)
			}
		}
	}

	if len(products) == 0 || (len(products) != lenIds) {
		return []Product{}, ErrRecordNotFound
	}

	return products, nil
}

func (p *ProductRepository) DeductStockProductById(id string, qtty int) error {
	items := p.Items
	lenItems := len(items)
	if lenItems == 0 {
		return ErrRecordNotFound
	}

	for i := 0; i < lenItems; i++ {
		if id == p.Items[i].Id {
			p.Items[i].Stock -= qtty
			return nil
		}
	}

	return ErrRecordNotFound
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	productRequest := ProductRequest{}

	date := time.Now().Format("2006-01-02")
	bytesBody, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		log.Printf("errReadAll: %v", errReadAll)
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errReadAll.Error(),
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if errUnmarshal := json.Unmarshal(bytesBody, &productRequest); errUnmarshal != nil {
		log.Printf("errUnmarshal: %v", errUnmarshal)
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errUnmarshal.Error(),
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productRequest.Name == "" {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Product name is empty.",
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productRequest.Stock == 0 {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Stock product is zero.",
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productRequest.Stock < 0 {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Stock product is minus.",
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productRequest.Price == 0 {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Price product is zero.",
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if productRequest.Price < 0 {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Price product is minus.",
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	productId, errGenerateProductId := generateProductId(len(productRepository.Items))
	if errGenerateProductId != nil {
		go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errGenerateProductId.Error(),
			date).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productRequest.id = productId
	productRepository.Add(productRequest)
	go monitoring.RequestTotalEndpointAddProduct.WithLabelValues(
		r.Method,
		strconv.Itoa(http.StatusOK),
		"Success",
		date).Inc()
	w.WriteHeader(http.StatusOK)
	return
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	monitoringDate := time.Now().Format("2006-01-02")
	if len(productRepository.Items) == 0 {
		go monitoring.RequestTotalEndpointGetProducts.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusNotFound),
			"Data empty.",
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	bytesItems, errMarshal := json.Marshal(&productRepository.Items)
	if errMarshal != nil {
		go monitoring.RequestTotalEndpointGetProducts.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errMarshal.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go monitoring.RequestTotalEndpointGetProducts.WithLabelValues(
		r.Method,
		strconv.Itoa(http.StatusOK),
		"Success",
		monitoringDate,
	).Inc()
	w.Write([]byte(bytesItems))
	return
}

func getProductsByIds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.Lock()
	defer mu.Unlock()

	monitoringDate := time.Now().Format("2006-01-02")
	if errParseForm := r.ParseForm(); errParseForm != nil {
		log.Printf("errParseForm: %v", errParseForm)
		go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errParseForm.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ids := r.Form["ids"]
	if len(ids) == 0 {
		go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusNotFound),
			"ids empty.",
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusNotFound)
		return
	}

	lenIds := len(ids)
	for i := 0; i < lenIds; i++ {
		if strings.Trim(ids[i], " ") == "" {
			go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
				r.Method,
				strconv.Itoa(http.StatusBadRequest),
				"id of ids empty.",
				monitoringDate,
			).Inc()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	products, errFindProductsByIds := productRepository.FindProductsByIds(ids)
	if errFindProductsByIds != nil {
		if errors.Is(errFindProductsByIds, ErrRecordNotFound) {
			go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
				r.Method,
				strconv.Itoa(http.StatusNotFound),
				errFindProductsByIds.Error(),
				monitoringDate,
			).Inc()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusInternalServerError),
			errFindProductsByIds.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytesProducts, errMarshal := json.Marshal(products)
	if errMarshal != nil {
		go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errMarshal.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go monitoring.RequestTotalEndpointGetProductsByIds.WithLabelValues(
		r.Method,
		strconv.Itoa(http.StatusOK),
		"Success",
		monitoringDate,
	).Inc()
	w.Write(bytesProducts)
	return
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	monitoringDate := time.Now().Format("2006-01-02")
	if id == "" {
		go monitoring.RequestTotalEndpointGetProductById.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			"Id empty.",
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	product, errFindProductById := productRepository.FindProductById(id)
	if errFindProductById != nil {
		if errors.Is(errFindProductById, ErrRecordNotFound) {
			go monitoring.RequestTotalEndpointGetProductById.WithLabelValues(
				r.Method,
				strconv.Itoa(http.StatusNotFound),
				errFindProductById.Error(),
				monitoringDate,
			).Inc()
			w.WriteHeader(http.StatusNotFound)
			return
		}

		go monitoring.RequestTotalEndpointGetProductById.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusInternalServerError),
			errFindProductById.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytesProduct, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		go monitoring.RequestTotalEndpointGetProductById.WithLabelValues(
			r.Method,
			strconv.Itoa(http.StatusBadRequest),
			errMarshal.Error(),
			monitoringDate,
		).Inc()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go monitoring.RequestTotalEndpointGetProductById.WithLabelValues(
		r.Method,
		strconv.Itoa(http.StatusOK),
		"Success",
		monitoringDate,
	).Inc()
	w.Write(bytesProduct)
	return
}

func deductStockProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/api/v1/products", addProduct)
	router.Get("/api/v1/products", getProducts)
	router.Get("/api/v1/products/{id}", getProductById)
	router.Get("/api/v1/products/ids", getProductsByIds)
	router.Get("/metrics", promhttp.Handler().ServeHTTP)

	utils.ReceiveMessage(constant.TOPIC_DEDUC_QTTY_PRODUCT_FOR_ORDER)

	log.Printf("Running product-service on :%s", port)
	log.Fatalf("Error when listen and serve: %v", http.ListenAndServe(fmt.Sprintf(":%s", port), router))
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type (
	OrdersRequest struct {
		Orders []Order `json:"orders"`
	}

	Order struct {
		orderDetailId string
		orderId       string
		ProductId     string `json:"productId"`
		Qtty          int    `json:"qtty"`
		Price         int    `json:"price"`
		totalPrice    int
	}

	OrderDetail struct {
		Id         string `json:"id"`
		OrderId    string `json:"orderId"`
		ProductId  string `json:"productId"`
		Qtty       int    `json:"qtty"`
		Price      int    `json:"price"`
		TotalPrice int    `json:"totalPrice"`
	}

	OrderRepository struct {
		Items []OrderDetail
	}

	ProductResponse struct {
		Id        string    `json:"id"`
		Name      string    `json:"name"`
		Stock     int       `json:"stock"`
		Price     int       `json:"price"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

var (
	mu                sync.Mutex
	port              = "8082"
	ErrRecordNotFound = errors.New("record not found")
	orderRepository   OrderRepository
	ordersIds         []string
	ordersDetailIds   []string

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
	}
)

func (o *OrderRepository) CreateOrder(orderRequest Order) error {
	o.Items = append(o.Items, OrderDetail{
		Id:         orderRequest.orderDetailId,
		OrderId:    orderRequest.orderId,
		ProductId:  orderRequest.ProductId,
		Qtty:       orderRequest.Qtty,
		Price:      orderRequest.Price,
		TotalPrice: orderRequest.totalPrice,
	})

	return nil
}

func (o *OrderRepository) CreateOrders(ordersRequest []Order) error {
	lenOrdersRequest := len(ordersRequest)
	for i := 0; i < lenOrdersRequest; i++ {
		orderRequest := ordersRequest[i]

		o.Items = append(o.Items, OrderDetail{
			Id:         orderRequest.orderDetailId,
			OrderId:    orderRequest.orderId,
			ProductId:  orderRequest.ProductId,
			Qtty:       orderRequest.Qtty,
			Price:      orderRequest.Price,
			TotalPrice: orderRequest.totalPrice,
		})
	}

	return nil
}

func createOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bytesRequest, errReadAll := io.ReadAll(r.Body)
	if errReadAll != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ordersRequest := OrdersRequest{}
	if errUnmarshal := json.Unmarshal(bytesRequest, &ordersRequest); errUnmarshal != nil {
		log.Printf("errUnmarshal: %v", errUnmarshal)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	orders := ordersRequest.Orders
	lenOrders := len(orders)
	if lenOrders == 0 {
		log.Printf("lenOrders: %v", lenOrders)
		mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productIds := []string{}
	for i := 0; i < lenOrders; i++ {
		productIds = append(productIds, orders[i].ProductId)
	}

	// ennsure item of productIds is not empty
	paramIds := url.Values{}
	for _, productId := range productIds {
		paramIds.Add("ids", productId)
		if strings.Trim(productId, " ") == "" {
			log.Printf("productId: %v", "empty")
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}
	}

	// do check products to product-service
	request, errNewProduct := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://product-service:8081/api/v1/products/ids?%s", paramIds.Encode()),
		nil)
	if errNewProduct != nil {
		log.Printf("errNewProduct: %v", errNewProduct)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}

	response, errDo := client.Do(request)
	if errDo != nil {
		log.Printf("errDo: %v", errDo)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("response.Status: %v, http.StausOK: %v", response.Status, http.StatusOK)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}

	productsResponse := make([]ProductResponse, 0)
	bytesBody, errReadAll := io.ReadAll(response.Body)
	if errReadAll != nil {
		log.Printf("errReadAll: %v", errReadAll)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}
	if errUnmarshal := json.Unmarshal(bytesBody, &productsResponse); errUnmarshal != nil {
		log.Printf("errUnmarshal: %v", errUnmarshal)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}
	mu.Unlock()

	mu.Lock()
	orderId, errGenerateOrderId := generateOrderId(len(ordersIds))
	if errGenerateOrderId != nil {
		log.Printf("errGenerateId: %v", errGenerateOrderId)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}
	ordersIds = append(ordersIds, orderId)

	mapProducts := map[string]ProductResponse{}
	for i := 0; i < lenOrders; i++ {
		productResponse := productsResponse[i]
		mapProducts[productResponse.Id] = productResponse
	}

	for i := 0; i < lenOrders; i++ {
		orders[i].orderId = orderId
		orderDetailId, errGenerateOrderDetail := generateOrderDetailId(len(ordersDetailIds))
		if errGenerateOrderDetail != nil {
			log.Printf("errGenerateOrderDetail: %v", errGenerateOrderDetail)
			// reversal
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}
		ordersDetailIds = append(ordersDetailIds, orderDetailId)
		orders[i].orderDetailId = orderDetailId

		product, ok := mapProducts[orders[i].ProductId]
		if !ok {
			log.Printf("ok: %v, productId: %v", ok, orders[i].ProductId)
			// reversal
			// go publish.reversal(orderId)
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}

		if product.Price != orders[i].Price {
			log.Printf("product.Price: %v, orderPrice: %v", product.Price, orders[i].Price)
			// reversal
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}

		if product.Stock < orders[i].Qtty {
			log.Printf("product.Stock: %v, orderQtty: %v", product.Stock, orders[i].Qtty)
			// reversal
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}

		product.Stock -= orders[i].Qtty
		mapProducts[product.Id] = product

		orders[i].totalPrice = product.Price * orders[i].Qtty
	}
	mu.Unlock()

	log.Print("orderdsIds")
	log.Println(ordersIds)
	log.Printf("ordersDetailIds")
	log.Println(ordersDetailIds)
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
	router.Use(middleware.Recoverer)

	router.Post("/api/v1/orders", createOrders)

	log.Printf("Running order-service on: %s", port)
	log.Fatalf("Error when listen and server: %v", http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func generateOrderId(currentSize int) (string, error) {
	if currentSize <= -1 {
		return "", errors.New("currentSize should greaer than -1")
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
	return fmt.Sprintf("ORD%s%d", strings.Repeat("0", totalDigitZero), tempCurrentSize), nil
}

func generateOrderDetailId(currentSize int) (string, error) {
	if currentSize <= -1 {
		return "", errors.New("currentSize should greaer than -1")
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
	return fmt.Sprintf("ODD%s%d", strings.Repeat("0", totalDigitZero), tempCurrentSize), nil
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
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
)

var (
	mu                sync.Mutex
	ErrRecordNotFound = errors.New("record not found")
	orderRepository   OrderRepository
	ordersIds         []string

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mu.Lock()
	orderId, errGenerateOrderId := generateOrderId(len(ordersIds))
	if errGenerateOrderId != nil {
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}
	mu.Unlock()

	mu.Lock()
	orders := ordersRequest.Orders
	lenOrders := len(orders)
	if lenOrders == 0 {
		mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	productIds := []string{}
	for i := 0; i < lenOrders; i++ {
		productIds = append(productIds, orders[i].ProductId)
	}

	// ennsure item of productIds is not empty
	for _, productId := range productIds {
		if strings.Trim(productId, " ") == "" {
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}
	}

	// do check products to product-service
	request, errNewProduct := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost/api/v1/products?ids=%v", productIds),
		nil)
	if errNewProduct != nil {
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}

	response, errDo := client.Do(request)
	if errDo != nil {
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}

	mu.Lock()
	for i := 0; i < lenOrders; i++ {
		orders[i].orderId = orderId
		// orders[i].
	}
	mu.Unlock()
}

func main() {

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

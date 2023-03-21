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
	"github.com/sinulingga23/microservices/order-service/constant"
	"github.com/sinulingga23/microservices/order-service/utils"
)

type (
	OrderRepository struct {
		Items           []OrderDetail
		OrdersIds       []string
		OrdersDetailIds []string
	}

	OrderService struct {
		orderRepository OrderRepository
	}
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
	mu sync.Mutex

	port               = "8082"
	hostProductService = "product-service:8082"

	ErrRecordNotFound = errors.New("record not found")

	client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
	}
)

func NewOrderRepository() OrderRepository {
	return OrderRepository{
		Items:           make([]OrderDetail, 0),
		OrdersIds:       make([]string, 0),
		OrdersDetailIds: make([]string, 0),
	}
}

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

func NewOrderService(orderRepository OrderRepository) OrderService {
	return OrderService{orderRepository: orderRepository}
}

func (service *OrderService) createOrders(w http.ResponseWriter, r *http.Request) {
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
		fmt.Sprintf("http://%s/api/v1/products/ids?%s", hostProductService, paramIds.Encode()),
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

	orderId, errGenerateOrderId := utils.GenerateOrderId(len(service.orderRepository.OrdersIds))
	if errGenerateOrderId != nil {
		log.Printf("errGenerateId: %v", errGenerateOrderId)
		w.WriteHeader(http.StatusBadRequest)
		mu.Unlock()
		return
	}
	service.orderRepository.OrdersIds = append(service.orderRepository.OrdersIds, orderId)

	mapProducts := map[string]ProductResponse{}
	for i := 0; i < lenOrders; i++ {
		productResponse := productsResponse[i]
		mapProducts[productResponse.Id] = productResponse
	}

	for i := 0; i < lenOrders; i++ {
		orders[i].orderId = orderId

		orderDetailId, errGenerateOrderDetail := utils.GenerateOrderDetailId(len(service.orderRepository.OrdersDetailIds))
		if errGenerateOrderDetail != nil {
			log.Printf("errGenerateOrderDetail: %v", errGenerateOrderDetail)
			// reversal
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}

		service.orderRepository.OrdersDetailIds = append(service.orderRepository.OrdersDetailIds, orderDetailId)
		orders[i].orderDetailId = orderDetailId

		productId := orders[i].ProductId
		product, ok := mapProducts[productId]
		if !ok {
			log.Printf("ok: %v, productId: %v", ok, productId)
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

		qtty := orders[i].Qtty
		if product.Stock < qtty {
			log.Printf("product.Stock: %v, orderQtty: %v", product.Stock, qtty)
			// reversal
			w.WriteHeader(http.StatusBadRequest)
			mu.Unlock()
			return
		}

		product.Stock -= orders[i].Qtty
		mapProducts[product.Id] = product

		go func(orderId, orderDetailId string, qtty int) {
			message := struct {
				OrderId       string `json:"orderId"`
				OrderDetailId string `json:"orderDetailId"`
				ProductId     string `json:"productId"`
				Qtty          int    `json:"qtty"`
			}{OrderId: orderId, OrderDetailId: orderDetailId, ProductId: productId, Qtty: qtty}

			bytesMessage, errMarshal := json.Marshal(&message)
			if errMarshal != nil {
				log.Printf("Error when marshal message of topic kafka: %v", errMarshal)
			}

			if errPublishMessage := utils.PublishMessage(constant.TOPIC_DEDUC_QTTY_PRODUCT_FOR_ORDER, bytesMessage); errPublishMessage != nil {
				log.Printf("Error when send message to topic: %v, Errors: %v", constant.TOPIC_DEDUC_QTTY_PRODUCT_FOR_ORDER, errPublishMessage)
			}
		}(orderId, orderDetailId, qtty)

		totalPrice := product.Price * orders[i].Qtty
		orders[i].totalPrice = totalPrice

		service.orderRepository.CreateOrder(Order{
			orderDetailId: orderDetailId,
			orderId:       orderId,
			totalPrice:    totalPrice,
			ProductId:     productId,
			Qtty:          qtty,
		})
	}

	mu.Unlock()
	w.WriteHeader(http.StatusOK)
	return
}

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("HOST_PRODUCT_SERVICE") != "" {
		hostProductService = os.Getenv("HOST_PRODUCT_SERVICE")
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	orderRepository := NewOrderRepository()
	orderService := NewOrderService(orderRepository)

	router.Post("/api/v1/orders", orderService.createOrders)

	log.Printf("Running order-service on: %s", port)
	log.Fatalf("Error when listen and server: %v", http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

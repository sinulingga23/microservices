package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestTotalEndpointAddProduct = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_add_product",
		Help: "Show the request total endpoint add product",
	}, []string{"http_method", "http_status", "log_message", "date"})

	RequestTotalEndpointGetProducts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_get_products",
		Help: "Show the request total endpoint get products",
	}, []string{"http_method", "http_status", "log_message", "date"})

	RequestTotalEndpointGetProductById = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_get_product_by_id",
		Help: "Show the request total endpoint get product by id",
	}, []string{"http_method", "http_status", "log_message", "date"})

	RequestTotalEndpointGetProductsByIds = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_get_products_by_ids",
		Help: "Show the request total endpoint get products by ids",
	}, []string{"http_method", "http_status", "log_message", "date"})
)

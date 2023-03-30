package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestTotalEndpointCreateOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_create_orders",
		Help: "Show the request total endpoint create orders",
	}, []string{"http_method", "http_status", "log_message", "date"})
)

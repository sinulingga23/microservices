package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestTotalEndpointAddProduct = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total_endpoint_add_product",
		Help: "Show the request total endpoint add product",
	}, []string{"status_code", "log_message", "date"})
)

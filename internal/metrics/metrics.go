package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP-запросов",
		},
		[]string{"method", "path"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Время ответа по ручкам",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	PVZCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pvz_created_total",
			Help: "Сколько ПВЗ было создано",
		},
	)

	ReceptionsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "receptions_created_total",
			Help: "Сколько приёмок было создано",
		},
	)

	ProductsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "products_created_total",
			Help: "Сколько товаров было добавлено",
		},
	)
)

func Register() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		PVZCreated,
		ReceptionsCreated,
		ProductsCreated,
	)
}

package analytics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil) // порт по умолчанию для метрик
}

var (
	RequestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Общее количество HTTP-запросов",
		},
	)

	// Гистограмма — для измерения времени отклика
	ResponseDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Время ответа HTTP в секундах",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Гейдж — для отслеживания количества активных пользователей
	ActiveUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Количество активных пользователей в системе",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(ResponseDuration)
	prometheus.MustRegister(ActiveUsers)
}

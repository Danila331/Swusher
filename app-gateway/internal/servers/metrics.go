package servers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil) // порт по умолчанию для метрик
}

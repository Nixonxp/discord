package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

const (
	appName   = "gateway_service"
	namespace = appName
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "grpc",
		Name:      appName + "_histogram_response_time_seconds",
		Help:      "Время ответа от сервера",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
	}, []string{"path", "is_error"})
)

// PrometheusMiddleware implements mux.MiddlewareFunc.
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)
		statusCode := lrw.Status()
		isError := strconv.FormatBool(statusCode != http.StatusOK)

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.String(), isError))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

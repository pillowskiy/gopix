package metric

import (
	"log"
	"strconv"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics interface {
	IncHits(status int, method, path string)
	ObserveResponseTime(status int, method, path string, observeTime float64)
}
type PrometheusMetrics struct {
	HitsTotal prometheus.Counter
	Hits      *prometheus.CounterVec
	Times     *prometheus.HistogramVec
}

func CreateMetrics(address string, name string) (Metrics, error) {
	metric := new(PrometheusMetrics)
	metric.HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: name + "_hits_total",
	})

	if err := prometheus.Register(metric.HitsTotal); err != nil {
		return nil, err
	}

	metric.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits",
		},
		[]string{"status", "method", "path"},
	)

	if err := prometheus.Register(metric.Hits); err != nil {
		return nil, err
	}

	metric.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "method", "path"},
	)

	if err := prometheus.Register(metric.Times); err != nil {
		return nil, err
	}

	if err := prometheus.Register(prometheus.NewBuildInfoCollector()); err != nil {
		return nil, err
	}

	go func() {
		router := echo.New()
		router.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		if err := router.Start(address); err != nil {
			log.Fatal(err)
		}
	}()

	return metric, nil
}

func (metric *PrometheusMetrics) IncHits(status int, method, path string) {
	metric.HitsTotal.Inc()
	metric.Hits.WithLabelValues(strconv.Itoa(status), method, path).Inc()
}

func (metric *PrometheusMetrics) ObserveResponseTime(status int, method, path string, observeTime float64) {
	metric.Times.WithLabelValues(strconv.Itoa(status), method, path).Observe(observeTime)
}

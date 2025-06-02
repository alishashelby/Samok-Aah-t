package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type QueryMetrics struct {
	Duration *prometheus.HistogramVec
	Executed *prometheus.CounterVec
	Errors   *prometheus.CounterVec
}

func NewQueryMetrics() *QueryMetrics {
	return &QueryMetrics{
		Duration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "query_in_seconds",
			Help:    "Duration of query in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"query"}),
		Executed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "queries_executed_total",
			Help: "Total number of queries executed",
		}, []string{"query"}),
		Errors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "query_errors_total",
			Help: "Total number of failed queries",
		}, []string{"query"}),
	}
}

func (qm *QueryMetrics) Middleware(queryName string, function func() error) error {
	qm.Executed.WithLabelValues(queryName).Inc()

	timer := prometheus.NewTimer(qm.Duration.WithLabelValues(queryName))
	defer timer.ObserveDuration()

	if err := function(); err != nil {
		qm.Errors.WithLabelValues(queryName).Inc()
		return err
	}

	return nil
}

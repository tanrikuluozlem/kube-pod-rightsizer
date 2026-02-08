package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	podsAnalyzed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rightsizer_pods_analyzed_total",
		Help: "Total number of pods analyzed",
	})

	recommendationsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rightsizer_recommendations_total",
		Help: "Total number of recommendations generated",
	})

	cpuSavingsPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rightsizer_cpu_savings_percent",
		Help: "CPU savings percentage per pod",
	}, []string{"namespace", "pod"})

	memorySavingsPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "rightsizer_memory_savings_percent",
		Help: "Memory savings percentage per pod",
	}, []string{"namespace", "pod"})
)

func init() {
	prometheus.MustRegister(podsAnalyzed)
	prometheus.MustRegister(recommendationsTotal)
	prometheus.MustRegister(cpuSavingsPercent)
	prometheus.MustRegister(memorySavingsPercent)
}

func Handler() http.Handler {
	return promhttp.Handler()
}

func IncPodsAnalyzed() {
	podsAnalyzed.Inc()
}

func IncRecommendations() {
	recommendationsTotal.Inc()
}

func SetCPUSavings(namespace, pod string, percent float64) {
	cpuSavingsPercent.WithLabelValues(namespace, pod).Set(percent)
}

func SetMemorySavings(namespace, pod string, percent float64) {
	memorySavingsPercent.WithLabelValues(namespace, pod).Set(percent)
}

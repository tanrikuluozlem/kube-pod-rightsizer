package metrics

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// PodMetrics holds resource usage for a single pod
type PodMetrics struct {
	Name      string
	Namespace string
	CPU       int64 // in millicores (50m = 50)
	Memory    int64 // in bytes
}

// Client wraps the metrics server client
type Client struct {
	metricsClient metricsv.Interface
}

// NewClient creates a metrics client
func NewClient(metricsClient metricsv.Interface) *Client {
	return &Client{
		metricsClient: metricsClient,
	}
}

// GetPodMetrics fetches metrics for all pods in a namespace
func (c *Client) GetPodMetrics(ctx context.Context, namespace string) ([]PodMetrics, error) {
	podMetricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	var result []PodMetrics
	for _, pm := range podMetricsList.Items {
		var totalCPU, totalMemory int64

		// sum all containers
		for _, container := range pm.Containers {
			totalCPU += container.Usage.Cpu().MilliValue()
			totalMemory += container.Usage.Memory().Value()
		}

		result = append(result, PodMetrics{
			Name:      pm.Name,
			Namespace: pm.Namespace,
			CPU:       totalCPU,
			Memory:    totalMemory,
		})
	}

	return result, nil
}

// GetAllPodMetrics fetches metrics across all namespaces
func (c *Client) GetAllPodMetrics(ctx context.Context) ([]PodMetrics, error) {
	return c.GetPodMetrics(ctx, "")
}

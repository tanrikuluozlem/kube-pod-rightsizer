package controller

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"

	"github.com/ozlemtanrikulu/kube-pod-rightsizer/internal/metrics"
	"github.com/ozlemtanrikulu/kube-pod-rightsizer/internal/recommender"
)

type Controller struct {
	k8sClient     kubernetes.Interface
	metricsClient *metrics.Client
	interval      time.Duration // Add necessary fields here
}

func NewController(k8sClient kubernetes.Interface, metricsClient *metrics.Client, interval time.Duration) *Controller {
	return &Controller{
		k8sClient:     k8sClient,
		metricsClient: metricsClient,
		interval:      interval,
	}
}

// Run starts the controller loop
func (c *Controller) Run(ctx context.Context) error {
	log.Println("starting controller...")

	c.analyze(ctx)

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("controller stopped")
			return nil
		case <-ticker.C:
			c.analyze(ctx)
		}
	}
}

// analyze scans all pods and generates recommendations
func (c *Controller) analyze(ctx context.Context) {
	log.Println("analyzing pods...")

	// get all pods
	pods, err := c.k8sClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("failed to list pods: %v", err)
		return
	}

	// get metrics for all pods
	podMetrics, err := c.metricsClient.GetAllPodMetrics(ctx)
	if err != nil {
		log.Printf("failed to get metrics: %v", err)
		return
	}

	// build metrics map for quick lookup
	metricsMap := make(map[string]metrics.PodMetrics)
	for _, pm := range podMetrics {
		key := pm.Namespace + "/" + pm.Name
		metricsMap[key] = pm
	}

	// analyze each pod
	for _, pod := range pods.Items {
		c.analyzePod(pod, metricsMap)
	}
}

// analyzePod checks a single pod
func (c *Controller) analyzePod(pod corev1.Pod, metricsMap map[string]metrics.PodMetrics) {
	// skip non-running pods
	if pod.Status.Phase != corev1.PodRunning {
		return
	}

	key := pod.Namespace + "/" + pod.Name
	pm, ok := metricsMap[key]
	if !ok {
		return // no metrics yet
	}

	// calculate total requests from pod spec
	var requestedCPU, requestedMemory int64
	for _, container := range pod.Spec.Containers {
		if cpu := container.Resources.Requests.Cpu(); cpu != nil {
			requestedCPU += cpu.MilliValue()
		}
		if mem := container.Resources.Requests.Memory(); mem != nil {
			requestedMemory += mem.Value()
		}
	}

	// run analysis
	result := recommender.Analyze(recommender.PodResources{
		Name:            pod.Name,
		Namespace:       pod.Namespace,
		UsedCPU:         pm.CPU,
		UsedMemory:      pm.Memory,
		RequestedCPU:    requestedCPU,
		RequestedMemory: requestedMemory,
	})

	// log recommendation if found
	if result != nil {
		log.Printf("RECOMMENDATION: %s/%s - CPU: %dm -> %dm (save %d%%), Memory: %dMi -> %dMi (save %d%%)",
			result.Namespace,
			result.PodName,
			result.CurrentCPU,
			result.RecommendedCPU,
			result.CPUSavings,
			result.CurrentMemory/(1024*1024),
			result.RecommendedMemory/(1024*1024),
			result.MemorySavings,
		)
	}
}

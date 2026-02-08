package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/ozlemtanrikulu/kube-pod-rightsizer/internal/controller"
	"github.com/ozlemtanrikulu/kube-pod-rightsizer/internal/metrics"
)

func main() {
	log.Println("starting kube-pod-rightsizer...")

	// create kubernetes config
	config, err := getKubeConfig()
	if err != nil {
		log.Fatalf("failed to get kubeconfig: %v", err)
	}

	// create kubernetes client
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create k8s client: %v", err)
	}

	// create metrics client
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create metrics client: %v", err)
	}

	// wrap metrics client
	mc := metrics.NewClient(metricsClient)

	// create controller
	ctrl := controller.NewController(k8sClient, mc, 5*time.Minute)

	// start health server
	go startHealthServer()

	// setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for SIGTERM/SIGINT
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Printf("received signal: %v, shutting down...", sig)
		cancel()
	}()

	// run controller
	if err := ctrl.Run(ctx); err != nil {
		log.Fatalf("controller error: %v", err)
	}
}

// getKubeConfig returns kubernetes config (in-cluster or from kubeconfig)
func getKubeConfig() (*rest.Config, error) {
	// try in-cluster config first (when running in kubernetes)
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Println("using in-cluster config")
		return config, nil
	}

	// fallback to kubeconfig (local development)
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}

	log.Printf("using kubeconfig: %s", kubeconfig)
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func startHealthServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.Handle("/metrics", metrics.Handler())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("starting health server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("health server error: %v", err)
	}
}

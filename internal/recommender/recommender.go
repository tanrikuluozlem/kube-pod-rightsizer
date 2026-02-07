package recommender

// PodResources holds current usage and requests for a pod
type PodResources struct {
	Name            string
	Namespace       string
	UsedCPU         int64 // millicores
	UsedMemory      int64 // bytes
	RequestedCPU    int64 // millicores
	RequestedMemory int64 // bytes
}

// Recommendation holds the analysis result
type Recommendation struct {
	PodName           string
	Namespace         string
	CurrentCPU        int64
	CurrentMemory     int64
	RecommendedCPU    int64
	RecommendedMemory int64
	CPUSavings        int // percentage
	MemorySavings     int // percentage
}

// Analyze checks if pod is over-provisioned and returns recommendation
func Analyze(pod PodResources) *Recommendation {
	if pod.RequestedCPU == 0 && pod.RequestedMemory == 0 {
		return nil
	}

	recommendedCPU := int64(float64(pod.UsedCPU) * 1.2)
	recommendedMemory := int64(float64(pod.UsedMemory) * 1.2)

	if recommendedCPU < 10 {
		recommendedCPU = 10
	}

	if recommendedMemory < 32*1024*1024 {
		recommendedMemory = 32 * 1024 * 1024
	}

	cpuSavings := 0
	if pod.RequestedCPU > 0 {
		cpuSavings = int((float64(pod.RequestedCPU-recommendedCPU) / float64(pod.RequestedCPU)) * 100)
	}

	memorySavings := 0
	if pod.RequestedMemory > 0 {
		memorySavings = int((float64(pod.RequestedMemory-recommendedMemory) / float64(pod.RequestedMemory)) * 100)
	}

	if cpuSavings < 10 && memorySavings < 10 {
		return nil
	}

	return &Recommendation{
		PodName:           pod.Name,
		Namespace:         pod.Namespace,
		CurrentCPU:        pod.RequestedCPU,
		CurrentMemory:     pod.RequestedMemory,
		RecommendedCPU:    recommendedCPU,
		RecommendedMemory: recommendedMemory,
		CPUSavings:        cpuSavings,
		MemorySavings:     memorySavings,
	}
}

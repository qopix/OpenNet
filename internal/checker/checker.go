package checker

import (
	"net/http"
	"time"

	"github.com/qopix/OpenNet/internal/config"
)

type Result struct {
	Name    string
	Address string
	OK      bool
	Latency time.Duration
}

func Check(endpoint config.Endpoint) Result {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(endpoint.Address)

	latency := time.Since(start)

	if err != nil {
		return Result{
			Name:    endpoint.Name,
			Address: endpoint.Address,
			OK:      false,
			Latency: latency,
		}
	}

	defer resp.Body.Close()

	return Result{
		Name:    endpoint.Name,
		Address: endpoint.Address,
		OK:      resp.StatusCode >= 200 && resp.StatusCode < 500,
		Latency: latency,
	}
}

func CheckAll(endpoints []config.Endpoint) []Result {
	results := make([]Result, 0, len(endpoints))

	for _, endpoint := range endpoints {
		results = append(results, Check(endpoint))
	}

	return results
}

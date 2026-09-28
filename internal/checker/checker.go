package checker

import (
	"fmt"
	"net/http"
	"time"
)

type Endpoint struct {
	Name string
	URL  string
}

var DefaultEndpoints = []Endpoint{
	{
		Name: "Yandex",
		URL:  "https://ya.ru",
	},
	{
		Name: "MAX",
		URL:  "https://max.ru",
	},
	{
		Name: "VK",
		URL:  "https://vk.com",
	},
	{
		Name: "RuStore",
		URL:  "https://rustore.ru",
	},
}

func Check(endpoint Endpoint) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(endpoint.URL)

	latency := time.Since(start)

	if err != nil {
		fmt.Printf(
			"%-10s [FAILED] %s (%dms)\n",
			endpoint.Name,
			endpoint.URL,
			latency.Milliseconds(),
		)
		return
	}

	resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		fmt.Printf(
			"%-10s [OK]     %s (%dms)\n",
			endpoint.Name,
			endpoint.URL,
			latency.Milliseconds(),
		)
		return
	}

	fmt.Printf(
		"%-10s [HTTP %d] %s (%dms)\n",
		endpoint.Name,
		resp.StatusCode,
		endpoint.URL,
		latency.Milliseconds(),
	)
}

func CheckDefault() {
	fmt.Println("OpenNet connectivity check")
	fmt.Println()

	for _, endpoint := range DefaultEndpoints {
		Check(endpoint)
	}
}

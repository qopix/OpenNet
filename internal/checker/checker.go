package checker

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/qopix/OpenNet/internal/key"
)

type Result struct {
	URL     string
	OK      bool
	Latency time.Duration
}

func Check(address string) Result {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(address)

	latency := time.Since(start)

	if err != nil {
		return Result{
			URL:     address,
			OK:      false,
			Latency: latency,
		}
	}

	resp.Body.Close()

	return Result{
		URL:     address,
		OK:      resp.StatusCode >= 200 && resp.StatusCode < 500,
		Latency: latency,
	}
}

func Connect(value string) error {
	host, clientID, _, err := key.Parse(value)

	if err != nil {
		return err
	}

	fmt.Println("OpenNet connection")
	fmt.Println()
	fmt.Println("Node:", host)
	fmt.Println("Client:", clientID)
	fmt.Println()
	fmt.Println("Checking node...")

	address := "https://" + host

	parsed, err := url.Parse(address)

	if err != nil || parsed.Host == "" {
		return fmt.Errorf("invalid node address")
	}

	result := Check(address)

	if !result.OK {
		return fmt.Errorf("node is unreachable")
	}

	fmt.Printf(
		"Node reachable (%dms)\n",
		result.Latency.Milliseconds(),
	)

	fmt.Println()
	fmt.Println("OpenNet handshake ready.")
	fmt.Println("Protocol version: 1")

	return nil
}

func NormalizeEndpoint(value string) string {
	value = strings.TrimSpace(value)

	if strings.HasPrefix(value, "http://") ||
		strings.HasPrefix(value, "https://") {
		return value
	}

	return "https://" + value
}

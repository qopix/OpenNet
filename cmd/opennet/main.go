package main

import (
	"fmt"
	"log"

	"github.com/qopix/OpenNet/internal/checker"
	"github.com/qopix/OpenNet/internal/config"
	"github.com/qopix/OpenNet/internal/web"
)

func main() {
	fmt.Println("OpenNet v0.5")
	fmt.Println("Network connectivity framework")
	fmt.Println()

	cfg := config.Default()

	fmt.Println("Checking endpoints...")
	fmt.Println()

	results := checker.CheckAll(cfg.Endpoints)

	for _, result := range results {
		status := "FAILED"

		if result.OK {
			status = "OK"
		}

		fmt.Printf(
			"%-20s [%s] %s (%dms)\n",
			result.Name,
			status,
			result.Address,
			result.Latency.Milliseconds(),
		)
	}

	fmt.Println()
	fmt.Println("Starting local web panel...")
	fmt.Println("http://127.0.0.1:8765")
	fmt.Println()

	server := web.New(
		"127.0.0.1:8765",
		"data/private.key",
	)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

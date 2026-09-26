package main

import (
	"fmt"
	"log"

	"github.com/qopix/OpenNet/internal/checker"
	"github.com/qopix/OpenNet/internal/config"
)

func main() {
	fmt.Println("OpenNet v0.3")
	fmt.Println("Network diagnostics")
	fmt.Println()

	cfg := config.Default()

	fmt.Println("Checking available endpoints...")
	fmt.Println()

	results := checker.CheckAll(cfg.Endpoints)

	for _, result := range results {
		status := "FAILED"

		if result.OK {
			status = "OK"
		}

		fmt.Printf(
			"%-20s [%s] %s\n",
			result.Name,
			status,
			result.Address,
		)
	}

	fmt.Println()
	fmt.Println("Done.")

	if len(results) == 0 {
		log.Println("No endpoints configured.")
	}
}

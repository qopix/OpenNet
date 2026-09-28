package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qopix/OpenNet/internal/config"
	"github.com/qopix/OpenNet/internal/socks5"
)

const configPath = "data/config.json"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println("OpenNet v1.0")
			return

		case "status":
			fmt.Println("OpenNet v1.0")
			fmt.Println("SOCKS5: 127.0.0.1:1080")
			return

		case "stop":
			fmt.Println("OpenNet v1.0 runs in foreground.")
			fmt.Println("Stop it with Ctrl+C.")
			return
		}
	}

	fmt.Println("=================================")
	fmt.Println(" OpenNet v1.0")
	fmt.Println("=================================")

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[OpenNet] Configuration loaded.")
	fmt.Println("[OpenNet] Starting transport...")
	fmt.Println("[OpenNet] Starting SOCKS5...")

	server := socks5.New(cfg)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("[OpenNet] SOCKS5 error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("[OpenNet] Shield: ON")
	fmt.Println("[OpenNet] SOCKS5: 127.0.0.1:1080")
	fmt.Println("OpenNet is running.")
	fmt.Println("Press Ctrl+C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	fmt.Println("\n[OpenNet] Stopped.")
}

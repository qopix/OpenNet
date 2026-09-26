package main

import (
	"fmt"
	"log"

	"github.com/qopix/OpenNet/internal/config"
	"github.com/qopix/OpenNet/internal/web"
)

func main() {
	cfg := config.Default()

	fmt.Println("OpenNet v0.1")
	fmt.Println()
	fmt.Println("Local control panel:")
	fmt.Printf("http://%s\n", cfg.PanelAddress)
	fmt.Println()

	server := web.New(cfg.PanelAddress)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
package main

import (
	"fmt"
	"log"

	"github.com/qopix/OpenNet/internal/web"
)

func main() {
	fmt.Println("OpenNet v0.1")
	fmt.Println("Control panel: http://127.0.0.1:8765")

	if err := web.Start(); err != nil {
		log.Fatal(err)
	}
}
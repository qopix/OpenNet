package main

import (
	"fmt"
	"os"

	"github.com/qopix/OpenNet/internal/checker"
	"github.com/qopix/OpenNet/internal/key"
	"github.com/qopix/OpenNet/internal/node"
)

func usage() {
	fmt.Println("OpenNet v0.7")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  opennet")
	fmt.Println("  opennet check")
	fmt.Println("  opennet key")
	fmt.Println("  opennet connect <key>")
	fmt.Println("  opennet node")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("OpenNet v0.7")
		fmt.Println("Open network connectivity framework")
		fmt.Println()
		usage()
		return
	}

	switch os.Args[1] {

	case "check":
		checker.CheckDefault()

	case "key":
		if err := key.CreateAndShow(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

	case "connect":
		if len(os.Args) != 3 {
			fmt.Println("Usage: opennet connect <key>")
			os.Exit(1)
		}

		if err := node.Connect(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "Connection failed:", err)
			os.Exit(1)
		}

	case "node":
		if err := node.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "Node error:", err)
			os.Exit(1)
		}

	default:
		usage()
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/qopix/OpenNet/internal/checker"
	"github.com/qopix/OpenNet/internal/key"
)

func usage() {
	fmt.Println("OpenNet v0.6")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  opennet")
	fmt.Println("  opennet key")
	fmt.Println("  opennet connect <key>")
	fmt.Println("  opennet check <url>")
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("OpenNet v0.6")
		fmt.Println("Open network connectivity framework")
		fmt.Println()
		usage()
		return
	}

	switch os.Args[1] {
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

		if err := checker.Connect(os.Args[2]); err != nil {
			fmt.Fprintln(os.Stderr, "Connection failed:", err)
			os.Exit(1)
		}

	case "check":
		if len(os.Args) != 3 {
			fmt.Println("Usage: opennet check <url>")
			os.Exit(1)
		}

		result := checker.Check(os.Args[2])

		if result.OK {
			fmt.Printf(
				"OK %s (%dms)\n",
				result.URL,
				result.Latency.Milliseconds(),
			)
			return
		}

		fmt.Printf(
			"FAILED %s (%dms)\n",
			result.URL,
			result.Latency.Milliseconds(),
		)

	default:
		usage()
		os.Exit(1)
	}
}

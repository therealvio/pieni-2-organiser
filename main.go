package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pieni-2-organiser/internal"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)

	err := internal.Handler(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

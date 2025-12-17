package main

import (
	"fmt"
	"os"

	"go-admin/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application startup failed: %v\n", err)
		os.Exit(1)
	}
}
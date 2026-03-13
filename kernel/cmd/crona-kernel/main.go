package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"crona/kernel/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Printf("crona-kernel: %v", err)
		os.Exit(1)
	}
}

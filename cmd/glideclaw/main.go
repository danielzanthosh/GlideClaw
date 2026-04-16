package main

import (
	"context"
	"log"
	"os"

	"glideclaw/internal/app"
)

func main() {
	ctx := context.Background()
	application, err := app.New()
	if err != nil {
		log.Fatalf("init failed: %v", err)
	}

	if err := application.Run(ctx, os.Args[1:]); err != nil {
		log.Fatalf("run failed: %v", err)
	}
}

// Copyright 2026 Daniel
// Licensed under the Apache License, Version 2.0

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

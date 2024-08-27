package main

import (
	"context"
	"log"

	"github.com/aaronland/go-broadcaster/app/broadcast"
)

func main() {

	ctx := context.Background()
	err := broadcast.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run broadcast application, %v", err)
	}
}

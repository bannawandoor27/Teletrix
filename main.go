package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/teletrix/internal/audio"
	"github.com/teletrix/internal/pipeline"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize audio processor
	processor, err := audio.NewProcessor(ctx, log)
	if err != nil {
		log.Fatalf("Failed to initialize audio processor: %v", err)
	}
	defer processor.Close()

	// Initialize pipeline
	pipe, err := pipeline.NewGStreamerPipeline(ctx, log)
	if err != nil {
		log.Fatalf("Failed to initialize pipeline: %v", err)
	}
	defer pipe.Close()

	// Setup signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start processing
	if err := processor.Start(); err != nil {
		log.Fatalf("Failed to start audio processing: %v", err)
	}

	// Start pipeline
	if err := pipe.Start(); err != nil {
		log.Fatalf("Failed to start pipeline: %v", err)
	}

	// Wait for shutdown signal
	<-sigCh
	log.Info("Shutting down...")
}

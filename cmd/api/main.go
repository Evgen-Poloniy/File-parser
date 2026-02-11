package main

import (
	"context"
	"errors"
	"file-parser/internal/config"
	"file-parser/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Loading server config
	serverConfig, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("error when loading server config: %v\n", err)
	}

	server := server.NewServer(&serverConfig)

	// Launching server in goroutine
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v\n", err)
		}
	}()

	// Wait signal of interrupt (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutting down server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v\n", err)
	}

	log.Println("server exited properly")
}

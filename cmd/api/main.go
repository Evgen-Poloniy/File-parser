package main

import (
	"context"
	"errors"
	"file-parser/internal/config"
	"file-parser/internal/handler"
	"file-parser/internal/repository"
	"file-parser/internal/server"
	"file-parser/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("error when loading env CONFIG_PATH")
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("error when loading config: %v", err)
	}

	repository := repository.NewRepository(nil)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	server := server.NewServer(&config.Server, handler.InitHandlers(&config.Logger))

	go func() {
		log.Printf("server is running on %s:%s\n", config.Server.Host, config.Server.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Server.TimeForGracefulShutdown*time.Second,
	)
	defer cancel()

	log.Println("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited properly")
}

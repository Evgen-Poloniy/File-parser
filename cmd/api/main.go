package main

import (
	"context"
	"errors"
	"file-parser/internal/config"
	"file-parser/internal/database"
	"file-parser/internal/handler"
	"file-parser/internal/logutil"
	"file-parser/internal/parser"
	"file-parser/internal/repository"
	"file-parser/internal/scanner"
	"file-parser/internal/server"
	"file-parser/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
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

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logutil.SetLevel(logger, &config.Logger)
	logutil.SetFormatter(logger, &config.Logger)

	db, err := database.NewPostgreSQLDatabase(&config.Database)
	if err != nil {
		logger.Fatalf("database initialize error: %v", err)
	}
	defer func() {
		db.Close()
		logger.Println("database connection has been closed")
	}()
	logger.Println("database connection established successful")

	repository := repository.NewRepository(db)

	queryService := service.NewQueryService(repository)
	parserService := service.NewParserService(repository)

	handler := handler.NewHandler(queryService)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	server := server.NewServer(&config.Server, handler.InitHandlers(&config.Logger))
	go func() {
		defer wg.Done()

		logger.Printf("server is running on %s:%s", config.Server.Host, config.Server.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("server error: %v", err)
			cancel()
		}
	}()

	jobs := make(chan string, 100)
	defer close(jobs)

	// Launch scanner at goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		scanner := scanner.NewScanner(&config.Parser, logger)

		scanner.Scan(ctx, jobs)
	}()

	// Launch of workers at goroutines
	for i := 0; i < config.Parser.CountOfWorkers; i++ {
		wg.Add(1)

		fileParser := parser.NewFileParser(parserService, &config.Parser, logger)

		go func() {
			defer wg.Done()

			fileParser.ParseFile(ctx, jobs)
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		cancel()
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		config.Server.TimeForGracefulShutdown*time.Second,
	)
	defer shutdownCancel()

	logger.Println("shutting down server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}

	wg.Wait()

	logger.Println("server exited properly")
}

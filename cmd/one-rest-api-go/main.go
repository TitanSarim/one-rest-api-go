package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/titanSarim/one-rest-api-go/internal/config"
	"github.com/titanSarim/one-rest-api-go/internal/http/handlers/student"
	"github.com/titanSarim/one-rest-api-go/internal/storage/sqlite"
)

func main() {
	// Load configuration settings from the config file
	cfg := config.MustLoad()

	// Database setup would go here (currently missing)
	storage, err := sqlite.New(cfg)
	if(err != nil){
		log.Fatalf("Error creating database: %v", err.Error())
	}
	slog.Info("Storage initialized", slog.String("Env", cfg.Env), slog.String("version", "1.0.0"))

	// Create a new router (multiplexer) to handle HTTP routes
	router := http.NewServeMux()

	// Define a simple route that responds with a welcome message
	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// Create an HTTP server with the loaded configuration
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr, // Set server address from config
		Handler: router,              // Attach the router to the server
	}

	// Print a message indicating the server has started
	slog.Info("Server started",slog.String("address",cfg.Addr))
	fmt.Printf("Server started on port: %s\n", cfg.HTTPServer.Addr)

	// Create a channel to listen for OS signals (graceful shutdown)
	done := make(chan os.Signal, 1)

	// Notify the `done` channel when an interrupt (Ctrl+C) or termination signal is received
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine to avoid blocking
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err) // If the server crashes, panic with the error
		}
	}()

	// Wait until an OS signal is received (blocking call)
	<-done

	slog.Info("Shutting down the server...")

	// Create a timeout context for graceful shutdown (5 seconds timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the context is canceled to free resources

	// Gracefully shutdown the server
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped")
}

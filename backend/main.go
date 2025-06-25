package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Promptzy/terminal-devtool/backend/api"
	"github.com/Promptzy/terminal-devtool/backend/middleware"
)

const (
	DefaultPort     = "8080"
	DefaultHost     = "localhost"
	ShutdownTimeout = 5 * time.Second
)

func main() {
	fmt.Println("📼 Terminal DevTool Go Backend")
	fmt.Println("✨ Version 0.1.0")

	// Set up structured logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get base directory for media files (default to current directory)
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	fmt.Printf("📁 Media base directory: %s\n", baseDir)

	// Create the API handler
	apiHandler := api.NewHandler(baseDir)

	// Create a new mux router and apply middleware
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/process", apiHandler.ProcessMedia)
	mux.HandleFunc("/api/compare", apiHandler.CompareMedia)
	mux.HandleFunc("/api/info", apiHandler.GetMediaInfo)
	mux.HandleFunc("/api/compress", apiHandler.CompressMedia)

	// Register health check endpoints
	mux.HandleFunc("/health", apiHandler.HealthCheck)

	// Backward compatibility for simple health check
	mux.HandleFunc("/api/health", apiHandler.HealthCheck)

	// Apply middleware
	var rootHandler http.Handler = mux
	rootHandler = middleware.Recovery(rootHandler)
	rootHandler = middleware.Logger(rootHandler)
	rootHandler = middleware.CORS(rootHandler)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Create server
	address := fmt.Sprintf("%s:%s", DefaultHost, port)
	server := &http.Server{
		Addr:         address,
		Handler:      rootHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Server starting on %s\n", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("👋 Server successfully shut down")
}

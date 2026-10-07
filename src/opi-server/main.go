// Package main implements the OPI API server for the developer sandbox.
// This server emulates OPI-standard infrastructure offloading APIs, allowing
// developers to test their clients against a realistic OPI endpoint
// without requiring actual SmartNIC or DPU hardware.
//
// The server exposes a REST/JSON API on port 8080 that mirrors OPI resource
// naming conventions (e.g., /v1/ports, /v1/pipelines).
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/opi-p4-sandbox/opi-server/handlers"
)

var (
	restPort = flag.Int("port", 8080, "REST API server port")
	logLevel = flag.String("log-level", "info", "Log level: debug, info, warn, error")
)

func main() {
	flag.Parse()

	// Override with environment variables if set
	if p := os.Getenv("OPI_PORT"); p != "" {
		fmt.Sscanf(p, "%d", restPort)
	}
	if l := os.Getenv("OPI_LOG_LEVEL"); l != "" {
		*logLevel = l
	}

	log.Printf("╔══════════════════════════════════════════════════╗")
	log.Printf("║        OPI API Server — Developer Sandbox        ║")
	log.Printf("╠══════════════════════════════════════════════════╣")
	log.Printf("║  REST API  : http://0.0.0.0:%-5d                ║", *restPort)
	log.Printf("║  Log Level : %-6s                               ║", *logLevel)
	log.Printf("║  Endpoints :                                     ║")
	log.Printf("║    GET  /healthz          - Health check          ║")
	log.Printf("║    GET  /v1/ports         - List ports            ║")
	log.Printf("║    POST /v1/ports         - Create port           ║")
	log.Printf("║    GET  /v1/ports/:id     - Get port              ║")
	log.Printf("║    DEL  /v1/ports/:id     - Delete port           ║")
	log.Printf("║    GET  /v1/pipelines     - List pipelines        ║")
	log.Printf("║    POST /v1/pipelines     - Create pipeline       ║")
	log.Printf("║    GET  /v1/pipelines/:id - Get pipeline          ║")
	log.Printf("║    DEL  /v1/pipelines/:id - Delete pipeline       ║")
	log.Printf("╚══════════════════════════════════════════════════╝")

	// Initialize services
	networkSvc := handlers.NewNetworkService()
	pipelineSvc := handlers.NewPipelineService()

	// Register REST routes
	mux := http.NewServeMux()
	handlers.RegisterRESTRoutes(mux, networkSvc, pipelineSvc)

	// Wrap with logging middleware
	handler := loggingMiddleware(mux)

	// Start HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *restPort),
		Handler: handler,
	}

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		server.Close()
	}()

	log.Printf("OPI API server listening on :%d", *restPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to serve: %v", err)
	}
	log.Printf("Server stopped.")
}

// loggingMiddleware logs every incoming HTTP request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

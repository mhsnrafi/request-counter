package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mhsnrafi/request-counter/counter"
)

// Run starts the HTTP server and listens for incoming requests.
func Run(cnt *counter.Counter) {
	http.HandleFunc("/", handleRequest(cnt))

	startServer()
}

func handleRequest(cnt *counter.Counter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count := cnt.Count()
		err := cnt.Save()
		if err != nil {
			panic(err)
		}

		response := fmt.Sprintf("Requests in the last 60 seconds: %d", count)
		_, err = w.Write([]byte(response))

		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

// startServer starts the HTTP server and listens for incoming requests.
func startServer() {
	server := &http.Server{Addr: ":8080"}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	waitForShutdownSignal()

	shutdownServer(server)
}

// waitForShutdownSignal waits for a shutdown signal (e.g., SIGINT or SIGTERM) to stop the server.
func waitForShutdownSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// shutdownServer gracefully shuts down the HTTP server, waiting for in-progress requests to complete.
func shutdownServer(server *http.Server) {
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

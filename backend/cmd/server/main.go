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

	"github.com/BhargavHirpara/devhealth/internal/api"
	"github.com/BhargavHirpara/devhealth/internal/scanner"
	"github.com/BhargavHirpara/devhealth/internal/store"
)

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "devhealth.db")
	githubToken := os.Getenv("GITHUB_TOKEN")

	if githubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Initialize store
	st, err := store.New(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}
	defer st.Close()

	// Initialize GitHub client and scanner
	ctx := context.Background()
	ghClient := api.NewGitHubClient(ctx, githubToken)
	sc := scanner.New(ghClient)

	// Initialize API server
	srv := api.New(st, sc)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("DevHealth API server starting on :%s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-done
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

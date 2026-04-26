package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jluisv16/hcm-go/internal/config"
	employeeapp "github.com/jluisv16/hcm-go/internal/employees/application"
	"github.com/jluisv16/hcm-go/internal/employees/infrastructure/memory"
	employeehttp "github.com/jluisv16/hcm-go/internal/employees/interfaces/http"
	"github.com/jluisv16/hcm-go/internal/http/handlers"
	"github.com/jluisv16/hcm-go/internal/http/router"
	"github.com/jluisv16/hcm-go/internal/version"
)

func main() {
	cfg := config.Load()
	startedAt := time.Now().UTC()

	healthHandler := handlers.NewHealthHandler(cfg.AppName, version.Value, startedAt)
	employeeRepository := memory.NewRepository(memory.SeedEmployees())
	employeeService := employeeapp.NewService(employeeRepository)
	employeeHandler := employeehttp.NewHandler(employeeService)
	engine := router.New(cfg, healthHandler, employeeHandler)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("starting %s on port %s in %s mode", cfg.AppName, cfg.HTTPPort, cfg.AppEnv)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}

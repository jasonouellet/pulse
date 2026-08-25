package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelchi"

	coreHTTP "pulse/internal/core/adapters/http"
	corePostgres "pulse/internal/core/adapters/postgres"
	"pulse/pkg/database"
	"pulse/pkg/observability"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Project PULSE (Modular Monolith)...")

	// 1. OTEL Observability Setup
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	otelCollector := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	shutdownTracer, err := observability.InitTracer(ctx, "pulse-monolith", otelCollector)
	if err != nil {
		slog.Warn("OpenTelemetry initialization skipped", "reason", err.Error())
	} else {
		defer func() { _ = shutdownTracer(context.Background()) }()
	}

	// 2. Database Connection & Auto-Migrations
	dbCfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres_secret"),
		DBName:   getEnv("DB_NAME", "pulse_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode)

	if err := database.RunMigrations("migrations/core", dbURL); err != nil {
		slog.Error("Failed to apply core migrations", "error", err)
	}

	dbPool, err := database.NewPostgresPool(ctx, dbCfg)
	if err != nil {
		slog.Error("Critical database connection failure", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// 3. HTTP Router & OTEL Middleware Setup
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(otelchi.Middleware("pulse-monolith-api"))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := dbPool.Ping(r.Context()); err != nil {
			http.Error(w, "Database Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK","service":"pulse-monolith"}`))
	})

	// 4. Register Core Module Routes
	coreRepo := corePostgres.NewUserRepository(dbPool)
	coreHandler := coreHTTP.NewUserHandler(coreRepo)
	coreHandler.RegisterRoutes(r)

	// 5. Graceful HTTP Shutdown Server
	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("PULSE Monolith HTTP listening on :%s", port))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server shutdown unexpectedly", "error", err)
		}
	case sig := <-shutdown:
		slog.Info("Shutting down PULSE monolith cleanly...", "signal", sig.String())
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = server.Shutdown(shutdownCtx)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

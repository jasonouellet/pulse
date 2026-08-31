package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/riandyrn/otelchi"

	coreHTTP "pulse/internal/core/adapters/http"
	corePostgres "pulse/internal/core/adapters/postgres"
	"pulse/pkg/database"
	"pulse/pkg/observability"
)

// @title           Project PULSE API
// @version         1.0
// @description     Modular backend API for managing sports clubs (PULSE OS).
// @termsOfService  http://swagger.io/terms/

// @contact.name    Project PULSE team
// @contact.url     https://github.com/jasonouellet/pulse

// @license.name  Source-Available Non-Commercial (ADR-005)

// @host      localhost:8080
// @BasePath  /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting Project PULSE (Modular Backend)...")

	// 1. OTEL Observability Setup
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	otelCollector := getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	shutdownTracer, err := observability.InitTracer(ctx, "pulse-backend", otelCollector)
	if err != nil {
		slog.Warn("OpenTelemetry initialization skipped", "reason", err.Error())
	} else {
		defer func() { _ = shutdownTracer(context.Background()) }()
	}

	// 1.b. Redis Setup
	// rdb, cleanupRedis := initRedis()
	// defer cleanupRedis()

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

	if err := database.RunAllMigrations(dbURL); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
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
	// The upstream nginx (controlled by us) overrides X-Real-IP via
	// `proxy_set_header X-Real-IP $remote_addr;` — therefore no value
	// provided by the client can pass. See GHSA-3fxj-6jh8-hvhx for
	// the context on the old, obsolete, and vulnerable RealIP middleware.
	r.Use(middleware.ClientIPFromHeader("X-Real-IP"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(otelchi.Middleware("pulse-backend-api"))

	// Prometheus metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	apiConfig := huma.DefaultConfig("Project PULSE API", "1.0.0")
	apiConfig.OpenAPIPath = "/openapi"
	apiConfig.DocsPath = "/docs"
	apiConfig.DocsRenderer = huma.DocsRendererSwaggerUI
	api := humachi.New(r, apiConfig)

	huma.Register(api, huma.Operation{
		OperationID:   "liveness-check",
		Method:        http.MethodGet,
		Path:          "/livez",
		Summary:       "Check process liveness",
		Description:   "Verifies that the backend process can serve requests.",
		Tags:          []string{"Health"},
		DefaultStatus: http.StatusOK,
	}, liveHandler)
	huma.Register(api, huma.Operation{OperationID: "readiness-check", Method: http.MethodGet, Path: "/readyz", Summary: "Check service readiness", Tags: []string{"Health"}, Errors: []int{http.StatusServiceUnavailable}}, readyHandler(dbPool))
	huma.Register(api, huma.Operation{OperationID: "health-check", Method: http.MethodGet, Path: "/healthz", Summary: "Check service health", Tags: []string{"Health"}, Errors: []int{http.StatusServiceUnavailable}}, readyHandler(dbPool))

	// 4. Register Core Module Routes
	coreRepo := corePostgres.NewUserRepository(dbPool)
	coreHandler := coreHTTP.NewUserHandler(coreRepo)
	coreHandler.RegisterRoutes(api)

	teamRepo := corePostgres.NewTeamRepository(dbPool)
	teamHandler := coreHTTP.NewTeamHandler(teamRepo)
	teamHandler.RegisterRoutes(api)

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
		slog.Info(fmt.Sprintf("PULSE Backend HTTP listening on :%s", port))
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
		slog.Info("Shutting down PULSE backend cleanly...", "signal", sig.String())
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		_ = server.Shutdown(shutdownCtx)
	}
}

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type healthOutput struct {
	Body healthResponse
}

func liveHandler(_ context.Context, _ *struct{}) (*healthOutput, error) {
	return &healthOutput{Body: healthResponse{Status: "OK", Service: "pulse-backend"}}, nil
}

func readyHandler(dbPool *pgxpool.Pool) func(context.Context, *struct{}) (*healthOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		if err := dbPool.Ping(ctx); err != nil {
			return nil, huma.Error503ServiceUnavailable("Database unavailable")
		}

		return &healthOutput{Body: healthResponse{Status: "OK", Service: "pulse-backend"}}, nil
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func initRedis() (*redis.Client, func()) {
	// If we are in local development mode without Docker, we launch miniredis
	if os.Getenv("ENV") == "development" || os.Getenv("REDIS_ADDR") == "" {
		mr, err := miniredis.Run()
		if err != nil {
			log.Fatalf("Failed to start local MiniRedis: %v", err)
		}

		log.Printf("🚀 MiniRedis (In-Memory) started locally on %s", mr.Addr())

		client := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})

		// Cleanup function to close properly during shutdown
		cleanup := func() {
			_ = client.Close()
			mr.Close()
		}
		return client, cleanup
	}

	// Otherwise (staging/prod), we connect to the real Redis
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	cleanup := func() { _ = client.Close() }
	return client, cleanup
}

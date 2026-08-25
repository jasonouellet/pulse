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

	"pulse/pkg/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Démarrage de Project PULSE (Monolithe Modulaire)...")

	dbCfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres_secret"),
		DBName:   getEnv("DB_NAME", "pulse_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	slog.Info("Connexion au pool PostgreSQL...")
	dbPool, err := database.NewPostgresPool(ctx, dbCfg)
	if err != nil {
		slog.Error("Échec critique lors de l'initialisation de la base de données", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	slog.Info("Connexion PostgreSQL établie avec succès !")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := dbPool.Ping(r.Context()); err != nil {
			http.Error(w, "Database Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK","service":"pulse-monolith"}`))
	})

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("Serveur HTTP à l'écoute sur le port :%s", port))
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Erreur critique du serveur HTTP", "error", err)
		}

	case sig := <-shutdown:
		slog.Info("Signal d'arrêt reçu, fermeture propre des ressources...", "signal", sig.String())

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Fermeture forcée du serveur HTTP", "error", err)
			if err := server.Close(); err != nil {
				slog.Error("Erreur lors de la fermeture directe des sockets", "error", err)
			}
		}
	}

	slog.Info("Arrêt complet et propre de Project PULSE.")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

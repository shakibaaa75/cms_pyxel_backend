package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"cms-backend/config"
	"cms-backend/database"
	"cms-backend/routes"
	"cms-backend/utils"
)

func main() {
	// Logger setup
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Load environment variables
	config.LoadEnv()

	// Background cleanup task
	go utils.PruneOldLimiters()

	// Connect database
	database.ConnectMongo()

	// 🔥 IMPORTANT: Render PORT fix
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local fallback
	}

	// Router
	handler := routes.Router()

	// Server config
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server starting",
		"port", port,
		"frontend", config.FRONTEND_URL,
		"admin", config.ADMIN_URL,
	)

	// Start server
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

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
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	config.LoadEnv()
	go utils.PruneOldLimiters()
	database.ConnectMongo()

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           routes.Router(),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server starting",
		"addr", srv.Addr,
		"frontend", config.FRONTEND_URL,
		"admin", config.ADMIN_URL,
	)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}

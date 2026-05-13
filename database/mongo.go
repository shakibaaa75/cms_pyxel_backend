package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"cms-backend/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoDB *mongo.Database

func ConnectMongo() {
	uri := config.GetEnv("MONGODB_URI", "")
	dbName := config.GetEnv("MONGODB_DB", "pyxel_cms")

	if uri == "" {
		slog.Warn("MONGODB_URI not set — CMS routes disabled")
		return
	}

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		slog.Error("mongo connect failed", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("mongo ping failed", "err", err)
		os.Exit(1)
	}

	MongoDB = client.Database(dbName)
	slog.Info("MongoDB connected", "database", dbName)

	EnsureIndexes()
}

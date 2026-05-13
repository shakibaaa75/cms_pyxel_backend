package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EnsureIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	type spec struct {
		col    string
		models []mongo.IndexModel
	}

	specs := []spec{
		{"projects", []mongo.IndexModel{
			{Keys: bson.D{{Key: "accessCode", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "magicToken", Value: 1}}},
		}},
		{"blogs", []mongo.IndexModel{
			{Keys: bson.D{{Key: "siteId", Value: 1}, {Key: "slug", Value: 1}}, Options: options.Index().SetUnique(true)},
			{Keys: bson.D{{Key: "siteId", Value: 1}, {Key: "published", Value: 1}, {Key: "createdAt", Value: -1}}},
			{Keys: bson.D{{Key: "published", Value: 1}, {Key: "createdAt", Value: -1}}}, // for public queries without siteId
		}},
	}

	for _, s := range specs {
		if _, err := MongoDB.Collection(s.col).Indexes().CreateMany(ctx, s.models); err != nil {
			slog.Warn("index creation failed", "collection", s.col, "err", err)
		} else {
			slog.Info("indexes ready", "collection", s.col)
		}
	}
}

package database

import (
	"banking-ledger/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnectivity struct {
	URI      string
	User     string
	Password string
	Database string
	MaxPool  uint64
}

func NewMongoConnectivity(cfg *config.Mongo) DBconnection {
	return &MongoConnectivity{
		URI:      cfg.URI,
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		MaxPool:  cfg.MaxPool,
	}
}
func (M *MongoConnectivity) Connection() (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build MongoDB connection URI
	uri := M.URI
	if M.User != "" && M.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s",
			M.User, M.Password, M.URI, M.Database)
	}

	clientOpts := options.Client().ApplyURI(uri)
	if M.MaxPool > 0 {
		clientOpts.SetMaxPoolSize(M.MaxPool)
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("unable to ping MongoDB: %w", err)
	}

	return client, nil
}

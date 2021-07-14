package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Because of one db is used decided to store name as const string
const dbName = "acg_db"

// Store represents DB abstraction
type Store struct {
	DatabaseURL string
	db          *mongo.Client
}

// NewStore return new Store object
func NewStore(dburl string) *Store {
	return &Store{
		DatabaseURL: dburl,
	}
}

// context for mongo db with 15 seconds timeout
var dbctx, _ = context.WithTimeout(context.Background(), 15*time.Second)

// Open just open new connection
func (s *Store) Open() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(s.DatabaseURL))
	if err != nil {
		return err
	}
	err = client.Connect(dbctx)
	if err != nil {
		return err
	}

	if err = client.Ping(dbctx, nil); err != nil {
		return err
	}

	s.db = client

	return nil
}

// Close just close the connection
func (s *Store) Close() {
	s.db.Disconnect(dbctx)
}

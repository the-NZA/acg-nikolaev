package mongostore

import (
	"context"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Because of one DB is used decided to store name as const string
const dbName = "acg_db"

// Store represents abstraction for MongoDB
type Store struct {
	db                  *mongo.Client
	postRepository      *PostRepository
	categoryRepository  *CategoryRepository
	materialsRepository *MaterialRepository
	matCatRepository    *MatCatRepository
	userRepository      *UserRepository
}

// context for mongo db with 15 seconds timeout
var dbctx, _ = context.WithTimeout(context.Background(), 15*time.Second)

// NewStore return new Store object or error
func NewStore(dbURL string) (store.Storer, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		return nil, err
	}
	err = client.Connect(dbctx)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(dbctx, nil); err != nil {
		return nil, err
	}

	return &Store{
		db: client,
	}, nil
}

// Close just aborts the connection
func (s *Store) Close() {
	s.db.Disconnect(dbctx)
}

// Implement Storer interface
func (s *Store) Posts() store.IPostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
	}

	return s.postRepository
}

func (s *Store) Categories() store.ICategoryRepository {
	if s.categoryRepository != nil {
		return s.categoryRepository
	}

	s.categoryRepository = &CategoryRepository{
		store:          s,
		collectionName: "categories",
	}

	return s.categoryRepository
}

func (s *Store) Materials() store.IMaterialRepository {
	if s.materialsRepository != nil {
		return s.materialsRepository
	}

	s.materialsRepository = &MaterialRepository{
		store: s,
	}

	return s.materialsRepository
}

func (s *Store) MatCategoies() store.IMatCategoryRepository {
	if s.matCatRepository != nil {
		return s.matCatRepository
	}

	s.matCatRepository = &MatCatRepository{
		store: s,
	}

	return s.matCatRepository
}

func (s *Store) Users() store.IUserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

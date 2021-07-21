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

// MongoStore represents abstraction for MongoDB
type MongoStore struct {
	db                  *mongo.Client
	postRepository      *PostRepository
	categoryRepository  *CategoryRepository
	materialsRepository *MaterialRepository
	matCatRepository    *MatCatRepository
	userRepository      *UserRepository
	serviceRepository   *ServiceRepository
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

	return &MongoStore{
		db: client,
	}, nil
}

// Close just aborts the connection
func (s *MongoStore) Close() {
	s.db.Disconnect(dbctx)
}

/*
 * Implement Storer interface
 */
func (s *MongoStore) Posts() store.IPostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store:          s,
		collectionName: "posts",
	}

	return s.postRepository
}

func (s *MongoStore) Categories() store.ICategoryRepository {
	if s.categoryRepository != nil {
		return s.categoryRepository
	}

	s.categoryRepository = &CategoryRepository{
		store:          s,
		collectionName: "categories",
	}

	return s.categoryRepository
}

func (s *MongoStore) Materials() store.IMaterialRepository {
	if s.materialsRepository != nil {
		return s.materialsRepository
	}

	s.materialsRepository = &MaterialRepository{
		store:          s,
		collectionName: "materials",
	}

	return s.materialsRepository
}

func (s *MongoStore) MatCategories() store.IMatCategoryRepository {
	if s.matCatRepository != nil {
		return s.matCatRepository
	}

	s.matCatRepository = &MatCatRepository{
		store:          s,
		collectionName: "matcategories",
	}

	return s.matCatRepository
}

func (s *MongoStore) Users() store.IUserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store:          s,
		collectionName: "users",
	}

	return s.userRepository
}

func (s *MongoStore) Services() store.IServiceRepository {
	if s.serviceRepository != nil {
		return s.serviceRepository
	}

	s.serviceRepository = &ServiceRepository{
		store:          s,
		collectionName: "services",
	}

	return s.serviceRepository
}

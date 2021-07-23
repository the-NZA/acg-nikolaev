package store

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IPostRepository defines interface for post repository
type IPostRepository interface {
	Create(*models.Post) error
	FindBySlug(string) (*models.Post, error)
	FindByURL(string) (*models.Post, error)
	FindAll(bson.M) ([]*models.Post, error)
	Delete(primitive.ObjectID) error
}

// ICategoryRepository defines interface for category repository
type ICategoryRepository interface {
	Create(*models.Category) error
	FindByID(primitive.ObjectID) (*models.Category, error)
	FindBySlug(string) (*models.Category, error)
	FindAll(bson.M) ([]*models.Category, error)
	Delete(primitive.ObjectID) error
}

// IMaterialRepository defines interface for material repository
type IMaterialRepository interface {
	Create(*models.Material) error
	FindByID(primitive.ObjectID) (*models.Material, error)
	FindBySlug(string) (*models.Material, error)
	FindAll(bson.M) ([]*models.Material, error)
	Delete(primitive.ObjectID) error
}

// IMatCategoryRepository defines interface for material category repository
type IMatCategoryRepository interface {
	Create(*models.MatCategory) error
	FindByID(primitive.ObjectID) (*models.MatCategory, error)
	FindBySlug(string) (*models.MatCategory, error)
	FindAll(bson.M) ([]*models.MatCategory, error)
	Delete(primitive.ObjectID) error
}

// IUserRepository defines interface for user repository
type IUserRepository interface {
	Create(*models.User) error
	// Find(string) (*models.User, error)
	Delete(primitive.ObjectID) error
	Login(string, string) error
}

// IServiceRepository defines interface for service repository
type IServiceRepository interface {
	Create(*models.Service) error
	FindByID(primitive.ObjectID) (*models.Service, error)
	FindBySlug(string) (*models.Service, error)
	Delete(primitive.ObjectID) error
	FindAll(filter bson.M) ([]*models.Service, error)
}

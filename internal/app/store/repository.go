package store

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IPostRepository defines interface for post repository
type IPostRepository interface {
	Create(*models.Post) error
	Find(string) (*models.Post, error)
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
	Find(string) (*models.Material, error)
	Delete(primitive.ObjectID) error
}

// IMatCategoryRepository defines interface for material category repository
type IMatCategoryRepository interface {
	Create(*models.MatCategory) error
	Find(string) (*models.MatCategory, error)
	Delete(primitive.ObjectID) error
}

// IUserRepository defines interface for user repository
type IUserRepository interface {
	Create(*models.User) error
	Find(string) (*models.User, error)
	Delete(primitive.ObjectID) error
}

package store

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IUserRepository defines interface for user repository
type IUserRepository interface {
	Create() error
	Find() error
	Delete() error
}

// IPostRepository defines interface for post repository
type IPostRepository interface {
	Create(*models.Post) error
	Find(string) (*models.Post, error)
	Delete(primitive.ObjectID) error
}

// ICategoryRepository defines interface for category repository
type ICategoryRepository interface {
	Create(*models.Category) error
	Find(string) (*models.Category, error)
	Delete(primitive.ObjectID) error
}

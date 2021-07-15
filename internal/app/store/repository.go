package store

import "github.com/the-NZA/acg-nikolaev/internal/app/models"

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
	Delete(*models.Post) error // Maybe just pass objectId??
}

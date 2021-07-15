package mongostore

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository implements IUserRepository
type UserRepository struct {
	store *Store
}

// Create save new post
func (p UserRepository) Create(*models.User) error {
	return nil
}

// Find lookup post by it slug
func (p UserRepository) Find(slug string) (*models.User, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p UserRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

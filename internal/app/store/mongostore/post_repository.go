package mongostore

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostRepository implements IPostRepository
type PostRepository struct {
	store *Store
}

// Create save new post
func (p PostRepository) Create(*models.Post) error {
	return nil
}

// Find lookup post by it slug
func (p PostRepository) Find(slug string) (*models.Post, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p PostRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

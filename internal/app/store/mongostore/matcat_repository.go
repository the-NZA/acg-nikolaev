package mongostore

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatCatRepository implements IMatCatRepository
type MatCatRepository struct {
	store *Store
}

// Create save new post
func (p MatCatRepository) Create(*models.MatCategory) error {
	return nil
}

// Find lookup post by it slug
func (p MatCatRepository) Find(slug string) (*models.MatCategory, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p MatCatRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

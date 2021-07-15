package mongostore

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryRepository implements ICategoryRepository
type CategoryRepository struct {
	store *Store
}

// Create new category
func (p CategoryRepository) Create(*models.Category) error {
	return nil
}

// Find category by it slug
func (p CategoryRepository) Find(slug string) (*models.Category, error) {
	return nil, nil
}

// Delete just marks category as deleted
func (p CategoryRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

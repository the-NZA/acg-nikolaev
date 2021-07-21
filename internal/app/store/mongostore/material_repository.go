package mongostore

import (
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MaterialRepository implements IMaterialRepository
type MaterialRepository struct {
	store *MongoStore
}

// Create save new material
func (p MaterialRepository) Create(*models.Material) error {
	return nil
}

// Find lookup material by it slug
func (p MaterialRepository) Find(slug string) (*models.Material, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p MaterialRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

package mongostore

import (
	"context"
	"fmt"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryRepository implements ICategoryRepository
type CategoryRepository struct {
	store          *Store
	collectionName string
}

// Create new category
func (c *CategoryRepository) Create(cat *models.Category) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := c.store.db.Database(dbName)
	col := db.Collection(c.collectionName)

	// _, err := col.InsertOne(ctx, cat)
	if _, err := col.InsertOne(ctx, cat); err != nil {
		return err
	}

	return nil
}

// Find category by it slug
func (c *CategoryRepository) Find(slug string) (*models.Category, error) {
	fmt.Println("repository find")

	return nil, nil
}

// Delete just marks category as deleted
func (c *CategoryRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

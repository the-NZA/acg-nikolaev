package mongostore

import (
	"context"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	if err := cat.Validate(); err != nil {
		return err
	}

	fcat, _ := c.FindBySlug(cat.Slug)
	if fcat != nil {
		return helpers.ErrCategoryAlreadyExist
	}

	db := c.store.db.Database(dbName)
	col := db.Collection(c.collectionName)

	if _, err := col.InsertOne(ctx, cat); err != nil {
		return err
	}

	return nil
}

func (c *CategoryRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.Category, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := c.store.db.Database(dbName)
	col := db.Collection(c.collectionName)
	res := col.FindOne(ctx, filter, opts...)

	cat := &models.Category{}

	err := res.Decode(cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

// Find category by it ID
func (c *CategoryRepository) FindByID(ID primitive.ObjectID) (*models.Category, error) {
	return c.findOne(bson.M{"_id": ID, "deleted": false})
}

// FindBySlug finds category by it slug
func (c *CategoryRepository) FindBySlug(slug string) (*models.Category, error) {
	return c.findOne(bson.M{"slug": slug, "deleted": false})
}

// FindAll return all categories
func (c *CategoryRepository) FindAll(filter bson.M) ([]*models.Category, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := c.store.db.Database(dbName)
	col := db.Collection(c.collectionName)

	res, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	cats := make([]*models.Category, 0)

	err = res.All(ctx, &cats)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

// Delete just marks category as deleted
func (c *CategoryRepository) Delete(deletedID primitive.ObjectID) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := c.store.db.Database(dbName)
	col := db.Collection(c.collectionName)

	_, err := col.UpdateByID(ctx, deletedID, bson.M{"$set": bson.M{"deleted": true}})
	if err != nil {
		return err
	}

	return nil
}

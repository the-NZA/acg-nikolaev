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

// MatCatRepository implements IMatCatRepository
type MatCatRepository struct {
	store          *Store
	collectionName string
}

// Create new material category
func (m MatCatRepository) Create(matcat *models.MatCategory) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := matcat.Validate(); err != nil {
		return err
	}

	fcat, _ := m.FindBySlug(matcat.Slug)
	if fcat != nil {
		return helpers.ErrCategoryAlreadyExist
	}

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	if _, err := col.InsertOne(ctx, matcat); err != nil {
		return err
	}

	return nil
}

func (m *MatCatRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.MatCategory, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)
	res := col.FindOne(ctx, filter, opts...)

	matcat := &models.MatCategory{}

	err := res.Decode(matcat)
	if err != nil {
		return nil, err
	}

	return matcat, nil
}

// Find material category by slug
func (m MatCatRepository) FindBySlug(slug string) (*models.MatCategory, error) {
	return m.findOne(bson.M{"slug": slug, "deleted": false})
}

// Find material category by it ID
func (m MatCatRepository) FindByID(ID primitive.ObjectID) (*models.MatCategory, error) {
	return m.findOne(bson.M{"_id": ID, "deleted": false})
}

// FindAll return all material repositories with specified filter
func (m MatCatRepository) FindAll(filter bson.M) ([]*models.MatCategory, error) {
	return nil, nil
}

func (m MatCatRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	_, err := col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Delete marks material category as deleted
func (m MatCatRepository) Delete(deletedID primitive.ObjectID) error {
	return m.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

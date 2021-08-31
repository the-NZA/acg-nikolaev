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

// MaterialRepository implements IMaterialRepository
type MaterialRepository struct {
	store          *MongoStore
	collectionName string
}

// Create save new material
func (m MaterialRepository) Create(material *models.Material) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := material.Validate(); err != nil {
		return err
	}

	fmaterial, _ := m.FindBySlug(material.Slug)
	if fmaterial != nil {
		return helpers.ErrMaterialAlreadyExist
	}

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	if _, err := col.InsertOne(ctx, material); err != nil {
		return err
	}

	return nil
}

func (m *MaterialRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.Material, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)
	res := col.FindOne(ctx, filter, opts...)

	material := &models.Material{}

	err := res.Decode(material)
	if err != nil {
		return nil, err
	}

	return material, nil

}

// FindBySlug lookup material by it slug
func (m MaterialRepository) FindBySlug(slug string) (*models.Material, error) {
	return m.findOne(bson.M{"slug": slug, "deleted": false})
}

// FindByID lookup material by ID
func (m MaterialRepository) FindByID(ID primitive.ObjectID) (*models.Material, error) {
	return m.findOne(bson.M{"_id": ID, "deleted": false})
}

// FindAll return all materials by filter parameter
func (m MaterialRepository) FindAll(filter bson.M) ([]*models.Material, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	res, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	materials := make([]*models.Material, 0)

	err = res.All(ctx, &materials)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

// Find return slice of material with filter and find options
func (m MaterialRepository) Find(filter bson.M, opts ...*options.FindOptions) ([]*models.Material, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	res, err := col.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	materials := make([]*models.Material, 0)

	err = res.All(ctx, &materials)
	if err != nil {
		return nil, err
	}

	return materials, nil
}

func (m MaterialRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
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

// Update recieve material, validate it and try to update it
func (m MaterialRepository) Update(updatedMaterial *models.Material) error {
	if err := updatedMaterial.Validate(); err != nil {
		return err
	}

	return m.updateOne(bson.M{"_id": updatedMaterial.ID}, bson.M{"$set": updatedMaterial})
}

// Delete marks post as deleted
func (m MaterialRepository) Delete(deletedID primitive.ObjectID) error {
	return m.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

// Count return number of materials that match filter with opts
func (m MaterialRepository) Count(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := m.store.db.Database(dbName)
	col := db.Collection(m.collectionName)

	return col.CountDocuments(ctx, filter, opts...)
}

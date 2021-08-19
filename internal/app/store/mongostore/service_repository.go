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

// ServiceRepository implements IServiceRepository
type ServiceRepository struct {
	store          *MongoStore
	collectionName string
}

// Create save new service
func (s ServiceRepository) Create(service *models.Service) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Validate(); err != nil {
		return err
	}

	fservice, _ := s.FindBySlug(service.Slug)
	if fservice != nil {
		return helpers.ErrServiceAlreadyExist
	}

	db := s.store.db.Database(dbName)
	col := db.Collection(s.collectionName)

	if _, err := col.InsertOne(ctx, service); err != nil {
		return err
	}

	return nil
}

func (s *ServiceRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.Service, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.store.db.Database(dbName)
	col := db.Collection(s.collectionName)
	res := col.FindOne(ctx, filter, opts...)
	// res := col.FindOne(ctx, bson.M{"_id": ID, "deleted": false})

	service := &models.Service{}

	err := res.Decode(service)
	if err != nil {
		return nil, err
	}

	return service, nil

}

// FindBySlug lookup service by it slug
func (s ServiceRepository) FindBySlug(slug string) (*models.Service, error) {
	return s.findOne(bson.M{"slug": slug, "deleted": false})
}

// FindByID lookup service by it id
func (s ServiceRepository) FindByID(ID primitive.ObjectID) (*models.Service, error) {
	return s.findOne(bson.M{"_id": ID, "deleted": false})
}

// FindAll return all services with specified filter
func (s ServiceRepository) FindAll(filter bson.M) ([]*models.Service, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.store.db.Database(dbName)
	col := db.Collection(s.collectionName)

	res, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	services := make([]*models.Service, 0)

	err = res.All(ctx, &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (s ServiceRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := s.store.db.Database(dbName)
	col := db.Collection(s.collectionName)

	_, err := col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Update validate updated service and try to update it in db
func (s ServiceRepository) Update(updatedService *models.Service) error {
	if err := updatedService.Validate(); err != nil {
		return err
	}

	return s.updateOne(bson.M{"_id": updatedService.ID}, bson.M{"$set": updatedService})
}

// Delete marks service as deleted
func (s ServiceRepository) Delete(deletedID primitive.ObjectID) error {
	return s.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

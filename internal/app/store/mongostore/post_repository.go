package mongostore

import (
	"context"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostRepository implements IPostRepository
type PostRepository struct {
	store          *MongoStore
	collectionName string
}

// Create save new post
func (p PostRepository) Create(post *models.Post) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := post.Validate(); err != nil {
		return err
	}

	fpost, _ := p.FindBySlug(post.Slug)
	if fpost != nil {
		return helpers.ErrPostAlreadyExist
	}

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	if _, err := col.InsertOne(ctx, post); err != nil {
		return err
	}

	return nil
}

func (p *PostRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.Post, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)
	res := col.FindOne(ctx, filter, opts...)
	// res := col.FindOne(ctx, bson.M{"_id": ID, "deleted": false})

	post := &models.Post{}

	err := res.Decode(post)
	if err != nil {
		return nil, err
	}

	return post, nil

}

// FindBySlug lookup post by it slug
func (p PostRepository) FindBySlug(slug string) (*models.Post, error) {
	return p.findOne(bson.M{"slug": slug, "deleted": false})
}

// FindByID lookup post by it ID
func (p PostRepository) FindByID(ID primitive.ObjectID) (*models.Post, error) {
	return p.findOne(bson.M{"_id": ID, "deleted": false})
}

// FindAll return all posts with specified filter
func (p PostRepository) FindAll(filter bson.M) ([]*models.Post, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	res, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	err = res.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Find return all posts with passed filter and find options
// This method is real projection to db find method
func (p PostRepository) Find(filter bson.M, opts ...*options.FindOptions) ([]*models.Post, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	res, err := col.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	err = res.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Aggregate gives opportunity to create more complex queries including 'joins' and etc
func (p PostRepository) Aggregate(pipeline mongo.Pipeline, opts ...*options.AggregateOptions) ([]*models.Post, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	res, err := col.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	err = res.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// Count return number of posts that match filter with opts
func (p PostRepository) Count(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	return col.CountDocuments(ctx, filter, opts...)
}

func (p PostRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	_, err := col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Update recieve post, validate it and try to update it
func (p PostRepository) Update(updatedPost *models.Post) error {
	if err := updatedPost.Validate(); err != nil {
		return err
	}

	return p.updateOne(bson.M{"_id": updatedPost.ID}, bson.M{"$set": updatedPost})
}

// Delete marks post as deleted
func (p PostRepository) Delete(deletedID primitive.ObjectID) error {
	return p.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

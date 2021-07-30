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

	fpost, _ := p.FindByURL(post.URL)
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

// FindByURL return all posts by it URL
func (p PostRepository) FindByURL(URL string) (*models.Post, error) {
	return p.findOne(bson.M{"url": URL, "deleted": false})
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

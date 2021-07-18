package mongostore

import (
	"context"
	"errors"
	"time"

	"github.com/the-NZA/acg-nikolaev/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostRepository implements IPostRepository
type PostRepository struct {
	store          *Store
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
		return errors.New("post already with exist")
	}

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	if _, err := col.InsertOne(ctx, post); err != nil {
		return err
	}

	return nil
}

// FindBySlug lookup post by it slug
func (p PostRepository) FindBySlug(slug string) (*models.Post, error) {
	return nil, nil
}

// FindAll return all posts with specified filter
func (p PostRepository) FindAll(filter bson.M) ([]*models.Post, error) {
	return nil, nil
}

// FindByURL return all posts by it URL
func (p PostRepository) FindByURL(URL string) (*models.Post, error) {
	return nil, nil
}

// Delete marks post as deleted
func (p PostRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

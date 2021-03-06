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

// PageRepository implements IPageRepository
type PageRepository struct {
	store          *MongoStore
	collectionName string
}

// Create save new page
func (p PageRepository) Create(page *models.Page) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := page.Validate(); err != nil {
		return err
	}

	fpage, _ := p.FindByURL(page.URL)
	if fpage != nil {
		return helpers.ErrPageAlreadyExist
	}

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	if _, err := col.InsertOne(ctx, page); err != nil {
		return err
	}

	return nil
}

func (p *PageRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.Page, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)
	res := col.FindOne(ctx, filter, opts...)
	// res := col.FindOne(ctx, bson.M{"_id": ID, "deleted": false})

	page := &models.Page{}

	err := res.Decode(page)
	if err != nil {
		return nil, err
	}

	return page, nil

}

// FindBySlug lookup page by it slug
func (p PageRepository) FindBySlug(slug string) (*models.Page, error) {
	return p.findOne(bson.M{"slug": slug})
}

// FindByURL return page by it URL
func (p PageRepository) FindByURL(URL string) (*models.Page, error) {
	return p.findOne(bson.M{"url": URL})
}

// FindByID return page by it ID
func (p PageRepository) FindByID(ID primitive.ObjectID) (*models.Page, error) {
	return p.findOne(bson.M{"_id": ID})
}

// FindAll return all pages with specified filter
func (p PageRepository) FindAll(filter bson.M) ([]*models.Page, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := p.store.db.Database(dbName)
	col := db.Collection(p.collectionName)

	res, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	pages := make([]*models.Page, 0)

	err = res.All(ctx, &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (p PageRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
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

// Update validate update page model and try to update it in db
func (p PageRepository) Update(updatedPage *models.Page) error {
	if err := updatedPage.Validate(); err != nil {
		return err
	}

	return p.updateOne(bson.M{"_id": updatedPage.ID}, bson.M{"$set": updatedPage})
}

// Delete marks page as deleted
func (p PageRepository) Delete(deletedID primitive.ObjectID) error {
	return p.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

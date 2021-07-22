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

// UserRepository implements IUserRepository
type UserRepository struct {
	store          *MongoStore
	collectionName string
}

// Create save new post
func (u UserRepository) Create(usr *models.User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := usr.Validate(); err != nil {
		return err
	}

	// If username already taken
	fusr, _ := u.FindByUsername(usr.Username)
	if fusr != nil {
		return helpers.ErrUserAlreadyExist
	}

	// If email already taken
	fusr, _ = u.FindByEmail(usr.Email)
	if fusr != nil {
		return helpers.ErrEmailAlreadyExist
	}

	if err := usr.BeforeSave(); err != nil {
		return err
	}

	db := u.store.db.Database(dbName)
	col := db.Collection(u.collectionName)

	if _, err := col.InsertOne(ctx, usr); err != nil {
		return err
	}

	return nil
}

func (u UserRepository) findOne(filter bson.M, opts ...*options.FindOneOptions) (*models.User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := u.store.db.Database(dbName)
	col := db.Collection(u.collectionName)
	res := col.FindOne(ctx, filter, opts...)
	// res := col.FindOne(ctx, bson.M{"_id": ID, "deleted": false})

	user := &models.User{}

	err := res.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil

}

// FindByUsername look up user by his username
func (u UserRepository) FindByUsername(username string) (*models.User, error) {
	return u.findOne(bson.M{"username": username})
}

// FindByEmail look up user by his email
func (u UserRepository) FindByEmail(email string) (*models.User, error) {
	return u.findOne(bson.M{"email": email})
}

// Delete marks post as deleted
func (p UserRepository) Delete(deletedID primitive.ObjectID) error {
	return nil
}

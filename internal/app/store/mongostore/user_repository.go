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
	return u.findOne(bson.M{"username": username, "deleted": false})
}

// FindByEmail look up user by his email
func (u UserRepository) FindByEmail(email string) (*models.User, error) {
	return u.findOne(bson.M{"email": email, "deleted": false})
}

func (u UserRepository) updateOne(filter bson.M, update bson.M, opts ...*options.UpdateOptions) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := u.store.db.Database(dbName)
	col := db.Collection(u.collectionName)

	_, err := col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

// Delete marks user as deleted
func (u UserRepository) Delete(deletedID primitive.ObjectID) error {
	return u.updateOne(bson.M{"_id": deletedID}, bson.M{"$set": bson.M{"deleted": true}})
}

func (u UserRepository) Login(username, password, secret string) (string, time.Time, error) {
	fusr, err := u.FindByUsername(username)
	if err != nil {
		return "", time.Time{}, err
	}

	fusr.Password = password

	token, expTime, err := fusr.DoLogin(secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expTime, nil
}

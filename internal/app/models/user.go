package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// hashCost for password hashing
const hashCost = 15

// User represets each user
type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	Username          string             `bson:"username" json:"username"`
	EncryptedPassword string             `bson:"pswd" json:"-"`
	Email             string             `bson:"email,omitempty" json:"email,omitempty"`
}

// Validate user struct
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&u.Username, validation.Required, validation.RuneLength(8, 0)),
		validation.Field(&u.EncryptedPassword, validation.Required),
		validation.Field(&u.Email, validation.Required, is.EmailFormat),
	)
}

// HashPassword generate password from input string
func (u *User) HashPassword(pass string) error {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), hashCost)
	if err != nil {
		return err
	}

	u.EncryptedPassword = string(hashBytes)
	return nil
}

// ComparePassword checks equality of given string and hashed passwords
func (u User) ComparePassword(p string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(p))
}

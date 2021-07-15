package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represets each user
type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"_id"`
	Username          string             `bson:"username" json:"username"`
	EncryptedPassword string             `bson:"pswd" json:"-"`
	Email             string             `bson:"email,omitempty" json:"email,omitempty"`
}

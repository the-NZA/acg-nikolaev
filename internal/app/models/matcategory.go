package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// MatCategory represent each materials category
type MatCategory struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Title   string             `bson:"title,omitempty" json:"title,omitempty"`
	Slug    string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Desc    string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Deleted bool               `bson:"deleted" json:"-"`
}

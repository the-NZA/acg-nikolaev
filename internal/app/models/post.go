package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is a structure for each post
type Post struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	Excerpt    string             `bson:"excerpt,omitempty" json:"excerpt,omitempty"`
	URL        string             `bson:"url,omitempty" json:"url,omitempty"`
	CategoryID primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Time       time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	TimeString string             `bson:"timestring,omitempty" json:"timestring,omitempty"`
	MetaDesc   string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	PostImg    string             `bson:"postimg,omitempty" json:"postimg,omitempty"`
	Deleted    bool               `bson:"deleted" json:"-"`
	// PageData    []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
}

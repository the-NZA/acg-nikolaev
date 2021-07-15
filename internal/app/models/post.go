package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is a structure for each posts
type Post struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Excerpt     string             `bson:"excerpt,omitempty" json:"excerpt,omitempty"`
	URL         string             `bson:"url,omitempty" json:"url,omitempty"`
	Category    string             `bson:"category,omitempty" json:"category,omitempty"`
	CategoryURL string             `bson:"categoryurl,omitempty" json:"categoryurl,omitempty"`
	Time        time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	TimeString  string             `bson:"timestring,omitempty" json:"timestring,omitempty"`
	MetaDesc    string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	// PageData    []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
	PostImg string `bson:"postimg,omitempty" json:"postimg,omitempty"`
	Deleted bool   `bson:"deleted" json:"-"`
}

// Category represents structure for each post category
type Category struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	Deleted  bool               `bson:"deleted" json:"-"`
	// Posts    []Post             `bson:"posts,omitempty" json:"posts,omitempty"`
}

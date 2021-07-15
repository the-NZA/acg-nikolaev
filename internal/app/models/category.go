package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category represents structure for each post category
type Category struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	Deleted  bool               `bson:"deleted" json:"-"`
}

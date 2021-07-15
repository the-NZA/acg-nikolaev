package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Page is basic model for each page
type Page struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	MetaDesc string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Slug     string             `bson:"slug,omitempty" json:"slug,omitempty"`
	PageData []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
}

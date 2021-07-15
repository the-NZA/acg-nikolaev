package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Service is a structure for representing each service
type Service struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Img      *ServiceImage      `bson:"img,omitempty" json:"img,omitempty"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	Desc     string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Deleted  bool               `bson:"deleted" json:"-"`
}

// ServiceImage represets basic structure of service card image
type ServiceImage struct {
	URL string `bson:"url,omitempty" json:"url,omitempty"`
	Alt string `bson:"alt,omitempty" json:"alt,omitempty"`
}

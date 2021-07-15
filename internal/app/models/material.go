package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Material represent structure of each material
type Material struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	CategoryID primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Desc       string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Time       time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	TimeString string             `bson:"timestring,omitempty" json:"timestring,omitempty"`
	FileLink   string             `bson:"filelink,omitempty" json:"filelink,omitempty"`
	Deleted    bool               `bson:"deleted" json:"-"`
}

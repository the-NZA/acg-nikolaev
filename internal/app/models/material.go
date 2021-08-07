package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Material represent structure of each material
type Material struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	Title         string             `bson:"title,omitempty" json:"title,omitempty"`
	MatCategoryID primitive.ObjectID `bson:"matcategory_id,omitempty" json:"matcategory_id,omitempty"`
	Slug          string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Desc          string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Time          time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	FileLink      string             `bson:"filelink,omitempty" json:"filelink,omitempty"`
	Deleted       bool               `bson:"deleted" json:"-"`
}

// MaterialShow represents material category with slice of materials for redreding in the browser
type MaterialShow struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title,omitempty" json:"title,omitempty"`
	Slug      string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Desc      string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Deleted   bool               `bson:"deleted" json:"-"`
	Materials []*Material        `bson:"materials"`
}

// TimeString return formated time string
func (m Material) TimeString() string {
	return m.Time.Format("02.01.2006")
}

// Validate material struct
func (m Material) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&m.MatCategoryID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&m.Title, validation.Required, validation.RuneLength(5, 55)),
		validation.Field(&m.Desc, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&m.Time, validation.Required),
		validation.Field(&m.FileLink, validation.Required),
		validation.Field(&m.Slug, validation.Required),
	)
}

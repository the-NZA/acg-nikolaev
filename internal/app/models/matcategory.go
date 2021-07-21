package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MatCategory represent each materials category
type MatCategory struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Title   string             `bson:"title,omitempty" json:"title,omitempty"`
	Slug    string             `bson:"slug,omitempty" json:"slug,omitempty"`
	Desc    string             `bson:"desc,omitempty" json:"desc,omitempty"`
	Deleted bool               `bson:"deleted" json:"-"`
}

func (m MatCategory) URL() string {
	return "matcategory" + m.Slug
}

func (m MatCategory) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&m.Title, validation.Required, validation.RuneLength(5, 35)),
		validation.Field(&m.Desc, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&m.Slug, validation.Required),
	)
}

package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category represents structure for each post category
type Category struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	Slug     string             `bson:"slug,omitempty" json:"slug,omitempty"`
	MetaDesc string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	Deleted  bool               `bson:"deleted" json:"-"`
}

// URL returns format url with format: "/category/category_slug"
func (c Category) URL() string {
	return "/category/" + c.Slug
}

// Validate check struct fields for correctness
func (c Category) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&c.Title, validation.Required, validation.RuneLength(5, 55)),
		validation.Field(&c.Subtitle, validation.Required, validation.RuneLength(25, 255)),
		validation.Field(&c.MetaDesc, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&c.Slug, validation.Required),
	)
}

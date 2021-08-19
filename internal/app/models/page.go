package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Page is basic model for each page
type Page struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Title    string             `bson:"title,omitempty" json:"title,omitempty"`
	Subtitle string             `bson:"subtitle,omitempty" json:"subtitle,omitempty"`
	MetaDesc string             `bson:"desc,omitempty" json:"desc,omitempty"`
	URL      string             `bson:"url,omitempty" json:"url,omitempty"`
	PageData []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
}

// Validate page struct
func (p Page) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&p.Title, validation.Required, validation.RuneLength(5, 55)),
		validation.Field(&p.Subtitle, validation.Required, validation.RuneLength(30, 255)),
		validation.Field(&p.MetaDesc, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&p.URL, validation.Required, validation.When(p.URL != "/", validation.RuneLength(5, 255)).Else(validation.Skip)),
		validation.Field(&p.PageData, validation.NilOrNotEmpty),
	)
}

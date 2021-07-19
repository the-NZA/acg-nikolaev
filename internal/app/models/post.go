package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is a structure for each post
type Post struct {
	ID         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `bson:"title,omitempty" json:"title,omitempty"`
	Snippet    string             `bson:"snippet,omitempty" json:"snippet,omitempty"`
	Slug       string             `bson:"slug,omitempty" json:"slug,omitempty"`
	URL        string             `bson:"url,omitempty" json:"url,omitempty"`
	CategoryID primitive.ObjectID `bson:"category_id,omitempty" json:"category_id,omitempty"`
	Time       time.Time          `bson:"time,omitempty" json:"time,omitempty"`
	MetaDesc   string             `bson:"metadesc,omitempty" json:"metadesc,omitempty"`
	PostImg    string             `bson:"postimg,omitempty" json:"postimg,omitempty"`
	PageData   []Block            `bson:"pagedata,omitempty" json:"pagedata,omitempty"`
	Deleted    bool               `bson:"deleted" json:"-"`
}

// TimeString return formated time string
func (p Post) TimeString() string {
	return p.Time.Format("02.01.2006")
}

// Validate check struct fields for correctness
func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&p.Title, validation.Required, validation.RuneLength(5, 35)),
		validation.Field(&p.Snippet, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&p.Slug, validation.Required, validation.RuneLength(5, 255)),
		validation.Field(&p.URL, validation.Required, validation.RuneLength(5, 255)),
		validation.Field(&p.CategoryID, validation.Required, validation.By(helpers.CheckObjectID)),
		validation.Field(&p.MetaDesc, validation.Required, validation.RuneLength(50, 255)),
		validation.Field(&p.Time, validation.Required),
		validation.Field(&p.PostImg, validation.Required),
		validation.Field(&p.PageData, validation.NilOrNotEmpty),
	)
}

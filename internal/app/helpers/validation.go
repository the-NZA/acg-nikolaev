package helpers

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidObjectID = errors.New("ObjectID must be valid")

func CheckObjectID(value interface{}) error {
	switch value.(type) {
	case primitive.ObjectID:
		if !primitive.IsValidObjectID(value.(primitive.ObjectID).Hex()) {
			return ErrInvalidObjectID
		}
	default:
		return ErrInvalidObjectID
	}

	return nil
}

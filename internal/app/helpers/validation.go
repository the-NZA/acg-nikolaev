package helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

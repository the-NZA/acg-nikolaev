package helpers

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errInvalidObjectID = errors.New("ObjectID must be valid")

func CheckObjectID(value interface{}) error {
	switch value.(type) {
	case primitive.ObjectID:
		if !primitive.IsValidObjectID(value.(primitive.ObjectID).Hex()) {
			return errInvalidObjectID
		}
	default:
		return errInvalidObjectID
	}

	// if !primitive.IsValidObjectID(value.(string)) {
	// 	return errInvalidObjectID
	// }

	return nil
}

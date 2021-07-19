package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCheckObjectID(t *testing.T) {
	testCases := []struct {
		name       string
		input      interface{}
		wantOutput error
	}{
		{
			name:       "Real ObjectID",
			input:      primitive.NewObjectID(),
			wantOutput: nil,
		},
		{
			name:       "Empty ObjectID",
			input:      "ObjectID(\"\")",
			wantOutput: ErrInvalidObjectID,
		},
		{
			name:       "Incorrect type (string)",
			input:      "ID(123)",
			wantOutput: ErrInvalidObjectID,
		},
		{
			name:       "Incorrect type (int)",
			input:      123,
			wantOutput: ErrInvalidObjectID,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			res := CheckObjectID(testCase.input)

			assert.Equal(t, testCase.wantOutput, res)
		})
	}
}

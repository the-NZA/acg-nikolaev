package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		wantOutput string
	}{
		{
			name:       "String 'Первая'",
			input:      "Первая",
			wantOutput: "pervaya",
		},
		{
			name:       "String 'ЛУЧШИЙ ОТрывок'",
			input:      "ЛУЧШИЙ ОТрывок",
			wantOutput: "luchshiy_otryvok",
		},
		{
			name:       "String '(empty)'",
			input:      "",
			wantOutput: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			res := GenerateSlug(testCase.input)

			assert.Equal(t, testCase.wantOutput, res)
		})
	}
}

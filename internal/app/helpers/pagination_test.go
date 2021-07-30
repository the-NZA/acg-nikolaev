package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePagination(t *testing.T) {
	testCases := []struct {
		name          string
		currentPage   uint
		numberOfPages uint
	}{
		{
			name:          "1",
			currentPage:   1,
			numberOfPages: 5,
		},
		{
			name:          "2",
			currentPage:   3,
			numberOfPages: 6,
		},
		{
			name:          "3",
			currentPage:   1,
			numberOfPages: 10,
		},
		{
			name:          "4",
			currentPage:   2,
			numberOfPages: 10,
		},
		{
			name:          "5",
			currentPage:   3,
			numberOfPages: 10,
		},
		{
			name:          "6",
			currentPage:   4,
			numberOfPages: 10,
		},
		{
			name:          "7",
			currentPage:   5,
			numberOfPages: 10,
		},
		{
			name:          "8",
			currentPage:   6,
			numberOfPages: 10,
		},
		{
			name:          "9",
			currentPage:   7,
			numberOfPages: 10,
		},
		{
			name:          "10",
			currentPage:   8,
			numberOfPages: 10,
		},
		{
			name:          "11",
			currentPage:   9,
			numberOfPages: 10,
		},
		{
			name:          "12",
			currentPage:   10,
			numberOfPages: 10,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			res := GeneratePagination(testCase.currentPage, testCase.numberOfPages)
			t.Logf("cur: %v, max: %v\n%v\n", testCase.currentPage, testCase.numberOfPages, res)

			assert.NotEmpty(t, res)

		})
	}
}

package helpers

import "strconv"

type PaginationLink struct {
	Link  string
	Value string
}

func GeneratePagination(currentPage, numberOfPages uint) []PaginationLink {
	pagination := make([]PaginationLink, 0)
	var num string

	// First link
	if currentPage != 1 {
		num = strconv.Itoa(int(currentPage - 1))
		pagination = append(pagination, PaginationLink{Link: num, Value: "←"})
	} else {
		pagination = append(pagination, PaginationLink{Link: "", Value: "←"})
	}

	// Last link
	if currentPage != numberOfPages {
		num = strconv.Itoa(int(currentPage + 1))
		pagination = append(pagination, PaginationLink{Link: num, Value: "→"})
	} else {
		pagination = append(pagination, PaginationLink{Link: "", Value: "→"})
	}
	return pagination
}

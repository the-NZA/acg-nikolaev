package helpers

import "strconv"

const (
	delta      = 2
	arrowLeft  = "←"
	arrowRight = "→"
)

type PaginationLink struct {
	Link  string
	Value string
}

func GeneratePagination(currentPage, numberOfPages uint) []PaginationLink {
	pagination := make([]PaginationLink, 0)
	var (
		num string
		i   uint

		dotsItem = PaginationLink{Link: "", Value: "..."}
	)

	// First link
	if currentPage != 1 {
		num = strconv.Itoa(int(currentPage - 1))
		pagination = append(pagination, PaginationLink{Link: num, Value: arrowLeft})
	} else {
		pagination = append(pagination, PaginationLink{Link: "", Value: arrowLeft})
	}

	if numberOfPages <= 6 {
		// Fill slice for 6 pages
		for i = uint(1); i <= numberOfPages; i++ {
			num = strconv.Itoa(int(i))
			pagination = append(pagination, PaginationLink{Link: num, Value: num})
		}
	} else {
		// Fill slice for a lot of pages with '...' for hidden pages
		switch currentPage {
		case 1:
			for i = currentPage; i <= currentPage+delta*2; i++ {
				num = strconv.Itoa(int(i))
				pagination = append(pagination, PaginationLink{Link: num, Value: num})
			}

			pagination = append(pagination, dotsItem)
		case 2:
			for i = currentPage - 1; i <= currentPage+delta; i++ {
				num = strconv.Itoa(int(i))
				pagination = append(pagination, PaginationLink{Link: num, Value: num})
			}

			pagination = append(pagination, dotsItem)
		case numberOfPages - 1:
			pagination = append(pagination, dotsItem)

			for i := currentPage - delta; i <= numberOfPages; i++ {
				num = strconv.Itoa(int(i))
				pagination = append(pagination, PaginationLink{Link: num, Value: num})
			}
		case numberOfPages:
			pagination = append(pagination, dotsItem)

			for i := currentPage - 2*delta; i <= numberOfPages; i++ {
				num = strconv.Itoa(int(i))
				pagination = append(pagination, PaginationLink{Link: num, Value: num})
			}
		default:
			if (currentPage - delta) != 1 {
				pagination = append(pagination, dotsItem)
			}

			for i := currentPage - delta; i <= currentPage+delta; i++ {
				num = strconv.Itoa(int(i))
				pagination = append(pagination, PaginationLink{Link: num, Value: num})
			}

			if (currentPage + delta) != numberOfPages {
				pagination = append(pagination, dotsItem)
			}
		}
	}

	// Last link
	if currentPage != numberOfPages {
		num = strconv.Itoa(int(currentPage + 1))
		pagination = append(pagination, PaginationLink{Link: num, Value: arrowRight})
	} else {
		pagination = append(pagination, PaginationLink{Link: "", Value: arrowRight})
	}

	return pagination
}

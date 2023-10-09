package helper

type PaginationResponse struct {
	TotalRecords int         `json:"total_records"`
	CurrentPage  int         `json:"current_page"`
	TotalPages   int         `json:"total_pages"`
	NextPage     interface{} `json:"next_page,omitempty"`
	PrevPage     interface{} `json:"prev_page,omitempty"`
}

func BuildPaginationResponse(totalRecords, page, pageSize int) PaginationResponse {
	totalPages := (totalRecords + pageSize - 1) / pageSize
	nextPage := 0
	if page < totalPages {
		nextPage = page + 1
	}
	prevPage := 0
	if page > 1 {
		prevPage = page - 1
	}

	return PaginationResponse{
		TotalRecords: totalRecords,
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     nextPage,
		PrevPage:     prevPage,
	}
}

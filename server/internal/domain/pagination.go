package domain

type Pagination[T any] struct {
	Items []T `json:"items"`
	PaginationInput
	Total int `json:"total"`
}

type PaginationInput struct {
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
}

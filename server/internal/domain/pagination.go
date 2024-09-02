package domain

type Pagination[T any] struct {
	Items   []T `json:"items"`
	Page    int `json:"page"`
	PerPage int `json:"perPage"`
	Total   int `json:"total"`
}

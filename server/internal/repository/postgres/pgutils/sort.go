package pgutils

import "fmt"

type sortOrder string

const (
	SortOrderASC  sortOrder = "ASC"
	SortOrderDESC sortOrder = "DESC"
)

type SortField struct {
	Field string
	Order sortOrder
}

type SortQueryBuilder struct {
	allowedFields map[string]SortField
}

func NewSortQueryBuilder() *SortQueryBuilder {
	return &SortQueryBuilder{
		allowedFields: make(map[string]SortField),
	}
}

func (s *SortQueryBuilder) AddField(name string, field SortField) *SortQueryBuilder {
	s.allowedFields[name] = field

	return s
}

func (s *SortQueryBuilder) GetSortField(name string) (SortField, bool) {
	field, ok := s.allowedFields[name]
	return field, ok
}

// Returns query string (field sortMethod) and true if sort field exists
func (s *SortQueryBuilder) SortQuery(name string) (string, bool) {
	field, ok := s.GetSortField(name)
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%s %s", field.Field, string(field.Order)), true
}

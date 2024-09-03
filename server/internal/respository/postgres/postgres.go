package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type InTxQueryCall func(tx *sqlx.Tx) (*sqlx.Rows, error)
type InTxQueryCallContext func(ctx context.Context, tx *sqlx.Tx) (*sqlx.Rows, error)

type sortOrder string

const (
	sortOrderASC  sortOrder = "ASC"
	sortOrderDESC sortOrder = "DESC"
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

func scanToStructSliceOf[T any](rows *sqlx.Rows) ([]T, error) {
	defer rows.Close()

	var dest []T
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}

		dest = append(dest, row)
	}

	return dest, nil
}

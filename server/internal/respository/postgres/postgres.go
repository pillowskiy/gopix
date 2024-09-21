package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	repository "github.com/pillowskiy/gopix/internal/respository"
	"github.com/pkg/errors"
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

	dest := []T{}
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, err
		}

		dest = append(dest, row)
	}

	return dest, nil
}

type Ext interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
}

type txKey struct{}
type PostgresRepository struct {
	db *sqlx.DB
}

func (r *PostgresRepository) ext(ctx context.Context) Ext {
	pgExt, ok := ctx.Value(txKey{}).(Ext)
	if !ok {
		return r.db
	}

	return pgExt
}

func (r *PostgresRepository) DoInTransaction(
	ctx context.Context, call repository.InTransactionalCall,
) (err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			xerr := tx.Rollback()
			if xerr != nil {
				err = errors.Wrap(err, xerr.Error())
			}

			fmt.Printf("Catched: %+v", err)
		} else {
			err = tx.Commit()
		}
	}()

	ctx = context.WithValue(ctx, txKey{}, tx)
	err = call(ctx)

	return
}

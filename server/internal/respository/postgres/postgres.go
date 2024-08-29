package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type InTxQueryCall func(tx *sqlx.Tx) (*sqlx.Rows, error)
type InTxQueryCallContext func(ctx context.Context, tx *sqlx.Tx) (*sqlx.Rows, error)

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

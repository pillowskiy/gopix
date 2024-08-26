package postgres

import (
	"github.com/jmoiron/sqlx"
)

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

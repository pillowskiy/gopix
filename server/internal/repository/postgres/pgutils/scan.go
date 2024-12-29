package pgutils

import "github.com/jmoiron/sqlx"

func ScanToStructSliceOf[T any](rows *sqlx.Rows) ([]T, error) {
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

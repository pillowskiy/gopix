package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pillowskiy/gopix/internal/repository"
	"github.com/pkg/errors"
)

type InTxQueryCall func(tx *sqlx.Tx) (*sqlx.Rows, error)
type InTxQueryCallContext func(ctx context.Context, tx *sqlx.Tx) (*sqlx.Rows, error)

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
	if tx := ctx.Value(txKey{}); tx != nil {
		return call(ctx)
	}

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

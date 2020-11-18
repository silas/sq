package sq

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Pool interface {
	Tx(ctx context.Context, fn func(tx Tx) error) error
}

type pgxPool struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Pool {
	return &pgxPool{pool: pool}
}

func (p *pgxPool) Tx(ctx context.Context, fn func(tx Tx) error) error {
	pgxtx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}
	tx := &Transaction{tx: pgxtx}
	return txExecute(ctx, tx, fn)
}

type Result pgconn.CommandTag

func (r Result) RowsAffected() int64 {
	return pgconn.CommandTag(r).RowsAffected()
}

type Tx interface {
	Exec(ctx context.Context, qb QueryBuilder) (Result, error)
	Query(ctx context.Context, qb QueryBuilder) (Rows, error)
	QueryRow(ctx context.Context, qb QueryBuilder) Row
}

type Transaction struct {
	tx pgx.Tx
}

func (tx *Transaction) Exec(ctx context.Context, qb QueryBuilder) (Result, error) {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	r, err := tx.tx.Exec(ctx, sql, args...)
	return Result(r), err
}

func (tx *Transaction) Query(ctx context.Context, qb QueryBuilder) (Rows, error) {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	return tx.tx.Query(ctx, sql, args...)
}

func (tx *Transaction) QueryRow(ctx context.Context, qb QueryBuilder) Row {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return rowError{err}
	}

	return tx.tx.QueryRow(ctx, sql, args...)
}

type rowError struct {
	err error
}

func (e rowError) Scan(...interface{}) error {
	return e.err
}

func IsError(err error, code string) bool {
	if err == nil {
		return false
	}

	var e *pgconn.PgError
	return errors.As(err, &e) && e != nil && e.Code == code
}

type Row interface {
	Scan(...interface{}) error
}

type Rows interface {
	Row
	Next() bool
	Close()
}

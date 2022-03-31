package sq

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrNoRows           = pgx.ErrNoRows
	ErrTxClosed         = pgx.ErrTxClosed
	ErrTxCommitRollback = pgx.ErrTxCommitRollback
)

type Config = pgxpool.Config

type Executor interface {
	Exec(ctx context.Context, qb StatementBuilder) (Result, error)
	Query(ctx context.Context, qb StatementBuilder) (Rows, error)
	QueryRow(ctx context.Context, qb StatementBuilder) Row
	All(ctx context.Context, qb StatementBuilder, dst interface{}) error
	One(ctx context.Context, qb StatementBuilder, dst interface{}) error
}

type Pool interface {
	Executor

	Tx(ctx context.Context, fn func(tx Tx) error) error
	Close()
}

func NewPool(p *pgxpool.Pool) Pool {
	return &pgxPool{pool: p}
}

func Connect(ctx context.Context, connString string) (Pool, error) {
	config, err := ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	return ConnectConfig(ctx, config)
}

func ParseConfig(connString string) (*Config, error) {
	return pgxpool.ParseConfig(connString)
}

func ConnectConfig(ctx context.Context, config *Config) (Pool, error) {
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return NewPool(pool), nil
}

type pgxPool struct {
	pool *pgxpool.Pool
}

func (p *pgxPool) Tx(ctx context.Context, fn func(tx Tx) error) error {
	pgxtx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return err
	}
	tx := &pgxTx{tx: pgxtx}
	return txExecute(ctx, tx, fn)
}

func (p *pgxPool) Close() {
	p.pool.Close()
}

func (p *pgxPool) Exec(ctx context.Context, qb StatementBuilder) (Result, error) {
	return exec(ctx, p.pool, qb)
}

func (p *pgxPool) Query(ctx context.Context, qb StatementBuilder) (Rows, error) {
	return query(ctx, p.pool, qb)
}

func (p *pgxPool) QueryRow(ctx context.Context, qb StatementBuilder) Row {
	return queryRow(ctx, p.pool, qb)
}

func (p *pgxPool) All(ctx context.Context, qb StatementBuilder, dst interface{}) error {
	return all(ctx, p.pool, qb, dst)
}

func (p *pgxPool) One(ctx context.Context, qb StatementBuilder, dst interface{}) error {
	return one(ctx, p.pool, qb, dst)
}

type Result = pgconn.CommandTag

type Tx interface {
	Executor
}

type pgxTx struct {
	tx pgx.Tx
}

func (tx *pgxTx) Exec(ctx context.Context, qb StatementBuilder) (Result, error) {
	return exec(ctx, tx.tx, qb)
}

func (tx *pgxTx) Query(ctx context.Context, qb StatementBuilder) (Rows, error) {
	return query(ctx, tx.tx, qb)
}

func (tx *pgxTx) QueryRow(ctx context.Context, qb StatementBuilder) Row {
	return queryRow(ctx, tx.tx, qb)
}

func (tx *pgxTx) All(ctx context.Context, qb StatementBuilder, dst interface{}) error {
	return all(ctx, tx.tx, qb, dst)
}

func (tx *pgxTx) One(ctx context.Context, qb StatementBuilder, dst interface{}) error {
	return one(ctx, tx.tx, qb, dst)
}

func exec(ctx context.Context, e pgxExecutor, qb StatementBuilder) (Result, error) {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	sql, err = replacePlaceholders(sql)
	if err != nil {
		return nil, err
	}

	return e.Exec(ctx, sql, args...)
}

func query(ctx context.Context, e pgxExecutor, qb StatementBuilder) (pgx.Rows, error) {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	sql, err = replacePlaceholders(sql)
	if err != nil {
		return nil, err
	}

	return e.Query(ctx, sql, args...)
}

func queryRow(ctx context.Context, e pgxExecutor, qb StatementBuilder) pgx.Row {
	sql, args, err := qb.ToSQL()
	if err != nil {
		return rowError{err}
	}

	sql, err = replacePlaceholders(sql)
	if err != nil {
		return rowError{err}
	}

	return e.QueryRow(ctx, sql, args...)
}

func all(ctx context.Context, e pgxExecutor, qb StatementBuilder, dst interface{}) error {
	rows, err := query(ctx, e, qb)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dst, rows)
}

func one(ctx context.Context, e pgxExecutor, qb StatementBuilder, dst interface{}) error {
	rows, err := query(ctx, e, qb)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dst, rows)
}

type pgxExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
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

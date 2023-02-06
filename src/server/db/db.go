package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool

	ErrNotFound = errors.New("no rows found")
)

func Connect(ctx context.Context, connString string) error {
	var err error

	pool, err = pgxpool.New(ctx, connString)

	return err
}

func Cleanup() {
	pool.Close()
}

func scanRows[T interface{ ScanFields() []any }](rows pgx.Rows, f func() T) ([]T, error) {
	var (
		out []T
		err error
	)

	for rows.Next() {
		var row = f()
		err = rows.Scan(row.ScanFields()...)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}

	if len(out) == 0 {
		return nil, ErrNotFound
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

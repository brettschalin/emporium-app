package db

import (
	"context"

	"marketplace-server/entity"

	"github.com/google/uuid"
)

const (
	userSelectFields = `id,created,updated,name,email,is_admin`
)

func CreateUser(ctx context.Context, u *entity.User) (uuid.UUID, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	query := `INSERT INTO users(name,email,is_admin) VALUES ($1, $2, $3) RETURNING id;`

	var (
		id uuid.UUID
	)

	row := conn.QueryRow(ctx, query, u.Name, u.Email, u.IsAdmin)

	if err = row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT ` + userSelectFields + ` FROM users WHERE id = $1 LIMIT 1;`

	row := conn.QueryRow(ctx, query, id.String())

	var user = new(entity.User)
	err = row.Scan(user.ScanFields()...)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(ctx context.Context, id uuid.UUID) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM users WHERE id = $1;`
	_, err = conn.Exec(ctx, query, id.String())
	if err != nil {
		return err
	}
	return nil
}

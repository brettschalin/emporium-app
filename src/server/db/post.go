package db

import (
	"context"
	"fmt"
	"os"

	"marketplace-server/entity"

	"github.com/google/uuid"
)

const (
	postSelectFields = `id,created,updated,created_by,expiration,contents,approved,in_mod_queue`
)

func GetPosts(ctx context.Context, inModQueue bool) ([]*entity.Post, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	const query = `SELECT ` + postSelectFields + ` FROM posts WHERE in_mod_queue = $1 AND expiration > now() ORDER BY created DESC;`

	rows, err := conn.Query(ctx, query, inModQueue)

	if err != nil {
		return nil, err
	}
	return scanRows(rows, func() *entity.Post { return new(entity.Post) })
}

func GetPost(ctx context.Context, id uuid.UUID) (*entity.Post, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT ` + postSelectFields + ` FROM posts WHERE id = $1 LIMIT 1;`

	row := conn.QueryRow(ctx, query, id.String())

	var post = new(entity.Post)
	err = row.Scan(post.ScanFields()...)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func CreatePost(ctx context.Context, p *entity.Post) (uuid.UUID, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	query := `INSERT INTO posts(created_by,expiration,root_parent,direct_parent,contents) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var (
		id uuid.UUID
	)

	row := conn.QueryRow(ctx, query, p.CreatedBy, p.Expiration, p.RootParent, p.DirectParent, p.Contents)

	if err = row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func UpdatePost(ctx context.Context, id uuid.UUID, updates *entity.UpdatePostRequest) (bool, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var (
		fieldNames = []string{}
		fields     = []any{}
	)

	if c := updates.Contents; c != nil {
		fieldNames = append(fieldNames, "contents")
		fields = append(fields, *c)
	}

	if len(fieldNames) == 0 {
		return false, nil
	}

	query := `UPDATE posts SET `
	for i := 0; i < len(fieldNames); i++ {
		if i != 0 {
			query += ", "
		}

		query += fmt.Sprintf(`%s = %d`, fieldNames[i], i+1)
	}
	query += fmt.Sprintf(" WHERE id = %d", len(fieldNames)+1)

	fields = append(fields, id.String())

	fmt.Fprintln(os.Stderr, "QUERY: "+query)

	res, err := conn.Exec(ctx, query, fields...)
	if err != nil {
		return false, err
	}

	return res.RowsAffected() == 1, nil
}

func ApprovePosts(ctx context.Context, approved bool, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `UPDATE posts SET approved = $1, in_mod_queue = false WHERE id = ANY($2);`

	_, err = conn.Exec(ctx, query, approved, ids)
	if err != nil {
		return err
	}

	return nil
}

func DeletePost(ctx context.Context, id uuid.UUID) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM posts WHERE id = $1;`
	_, err = conn.Exec(ctx, query, id.String())
	if err != nil {
		return err
	}

	return nil
}

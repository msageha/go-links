package user

import (
	"context"
	"database/sql"

	db "github.com/mzk622/go-links/backend/pkg/db"
	errors "github.com/pkg/errors"
)

type User struct {
	ID    uint64
	Email string
}

const (
	SELECT_USER_BY_ID    = "SELECT id, email FROM user WHERE id = ? AND is_deleted = FALSE"
	SELECT_USER_BY_EMAIL = "SELECT id, email FROM user WHERE email = ? AND is_deleted = FALSE"
)

func FindUserByID(ctx context.Context, id uint64) (*User, error) {
	query := SELECT_USER_BY_ID
	result, err := db.Query(ctx, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to exec query %s", query)
	}
	defer func() { _ = result.Close() }()
	return toUser(result.Rows)
}

func FindUserByEmail(ctx context.Context, email string) (*User, error) {
	query := SELECT_USER_BY_EMAIL
	result, err := db.Query(ctx, query, email)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to exec query %s", query)
	}
	defer func() { _ = result.Close() }()
	return toUser(result.Rows)
}

func toUser(rows *sql.Rows) (*User, error) {
	if !rows.Next() {
		return nil, nil
	}

	var user User
	err := rows.Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to scan query")
	}

	return &user, nil
}

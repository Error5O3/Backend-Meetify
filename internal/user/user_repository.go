package user

import (
	"context"
	"database/sql"
)

// initiliaze database
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastID int64
	query := "INSERT INTO users(username, email) VALUES ($1, $2) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email).Scan(&lastID)
	if err != nil {
		return &User{}, err
	}
	user.ID = lastID
	return user, nil
}

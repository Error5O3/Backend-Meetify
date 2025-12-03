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

func (r *repository) LoginUser(ctx context.Context, user *User) (*User, error) {
	// Simple existence check: return nil if a user with the given username exists
	query := "SELECT id, username, email FROM users WHERE username = $1"
	var u User
	err := r.db.QueryRowContext(ctx, query, user.Username).Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		// sql.ErrNoRows means user not found; propagate so caller can handle
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetAvail(ctx context.Context, userID int64) (*UserGridResponse, error) {
	// Simple existence check: return nil if a user with the given username exists
	query := "SELECT time_slot_id FROM user_availability WHERE user_id = $1"
	rows, err := r.db.QueryContext(ctx, query, userID)

	if err != nil {
		// sql.ErrNoRows means user not found; propagate so caller can handle
		return nil, err
	}

	defer rows.Close()

	var slots []int64

	for rows.Next() {
		var slot int64
		err := rows.Scan(&slot)
		if err != nil {
			return nil, err
		}

		slots = append(slots, slot)
	}
	return &UserGridResponse{
		AvailSlots: slots,
	}, nil
}

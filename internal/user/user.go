package user

import "context"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Repository interface {
	CreateUser(c context.Context, user *User) (*User, error)
}
type Service interface {
	CreateUser(c context.Context, req *UserRequest) (*UserResponse, error)
}

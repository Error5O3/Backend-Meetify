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

type UserGridResponse struct {
	AvailSlots []int64 `json:"avail"`
}

type Repository interface {
	CreateUser(c context.Context, user *User) (*User, error)
	LoginUser(c context.Context, user *User) (*User, error)
	GetAvail(c context.Context, userID int64) (*UserGridResponse, error)
}

type Service interface {
	CreateUser(c context.Context, req *UserRequest) (*UserResponse, error)
	LoginUser(c context.Context, req *UserRequest) (*UserResponse, error)
	GetAvail(c context.Context, userID int64) (*UserGridResponse, error)
}

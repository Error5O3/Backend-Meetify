package user

import (
	"context"
	"time"
)

type service struct {
	Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo,
	}
}

func (s *service) CreateUser(c context.Context, req *UserRequest) (*UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	usr := &User{
		Username: req.Username,
		Email:    req.Email,
	}
	user, err := s.Repository.CreateUser(ctx, usr)
	if err != nil {
		return nil, err
	}
	res := &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	return res, nil
}

func(s *service) LoginUser(c context.Context, req *UserRequest) (*UserResponse , error){
	ctx, cancel := context.WithTimeout(c , 5*time.Second)
	defer cancel()
	usr := &User{
		Username: req.Username,
	}
	user, err := s.Repository.LoginUser(ctx, usr)
	if err != nil{
		return nil, err
	}
	res := &UserResponse{
		ID: user.ID,
		Username : user.Username,
		Email: user.Email,
	}
	return res,nil
}

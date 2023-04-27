package service

import (
	"context"
	"github.com/praslar/lib/common"
	"movieon_be/pkg/model"
	"movieon_be/pkg/repo"
)

type UserService struct {
	repo repo.PGInterface
}

type UserInterface interface {
	Create(ctx context.Context, ob model.UserRequest) (rs *model.User, err error)
	Update(ctx context.Context, ob model.UserRequest) (rs *model.User, err error)
	Delete(ctx context.Context, id string) (err error)
	GetOne(ctx context.Context, id string) (rs *model.User, err error)
	GetList(ctx context.Context, req model.UserParams) (rs *model.UserResponse, err error)
	Login(ctx context.Context, req model.LoginRequest) (rs *model.User, err error)
	Register(ctx context.Context, req model.RegisterRequest) (err error)
}

func NewUserService(repo repo.PGInterface) UserInterface {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req model.UserRequest) (rs *model.User, err error) {

	ob := &model.User{}
	common.Sync(req, ob)

	if err := s.repo.CreateUser(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *UserService) Update(ctx context.Context, req model.UserRequest) (rs *model.User, err error) {

	ob := &model.User{}
	common.Sync(req, ob)

	if err := s.repo.UpdateUser(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *UserService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) GetOne(ctx context.Context, id string) (rs *model.User, err error) {

	ob, err := s.repo.GetOneUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *UserService) GetList(ctx context.Context, req model.UserParams) (rs *model.UserResponse, err error) {

	ob, err := s.repo.GetListUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

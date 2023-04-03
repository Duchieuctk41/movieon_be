package service

import (
	"context"
	"github.com/praslar/lib/common"
	"movieon_be/pkg/model"
	"movieon_be/pkg/repo"
	"movieon_be/pkg/valid"
)

type MovieService struct {
	repo repo.PGInterface
}

type MovieInterface interface {
	Create(ctx context.Context, ob model.MovieRequest) (rs *model.Movie, err error)
	Update(ctx context.Context, ob model.MovieRequest) (rs *model.Movie, err error)
	Delete(ctx context.Context, id string) (err error)
	GetOne(ctx context.Context, id string) (rs *model.Movie, err error)
	GetList(ctx context.Context, req model.MovieParams) (rs *model.MovieResponse, err error)
}

func NewMovieService(repo repo.PGInterface) MovieInterface {
	return &MovieService{repo: repo}
}

func (s *MovieService) Create(ctx context.Context, req model.MovieRequest) (rs *model.Movie, err error) {

	ob := &model.Movie{}
	common.Sync(req, ob)

	// check poster exist. create default image when don't have poster
	if ob.Poster == "" {
		ob.Poster = "https://www.publicdomainpictures.net/pictures/280000/velka/not-found-image-15383864787lu.jpg"
	}

	if err := s.repo.CreateMovie(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *MovieService) Update(ctx context.Context, req model.MovieRequest) (rs *model.Movie, err error) {
	ob, err := s.repo.GetOneMovie(ctx, valid.String(req.ID))
	if err != nil {
		return nil, err
	}

	common.Sync(req, ob)

	if err := s.repo.UpdateMovie(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *MovieService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteMovie(ctx, id)
}

func (s *MovieService) GetOne(ctx context.Context, id string) (rs *model.Movie, err error) {

	ob, err := s.repo.GetOneMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *MovieService) GetList(ctx context.Context, req model.MovieParams) (rs *model.MovieResponse, err error) {

	ob, err := s.repo.GetListMovie(ctx, req)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

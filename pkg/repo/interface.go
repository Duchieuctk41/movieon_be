package repo

import (
	"context"
	"gorm.io/gorm"
	"movieon_be/pkg/model"
)

type PGInterface interface {
	GetRepo() *gorm.DB
	Transaction(ctx context.Context, f func(rp PGInterface) error) error
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)

	// movie
	CreateMovie(ctx context.Context, user *model.Movie) error
	UpdateMovie(ctx context.Context, user *model.Movie) error
	DeleteMovie(ctx context.Context, id string) error
	GetOneMovie(ctx context.Context, id string) (*model.Movie, error)
	GetListMovie(ctx context.Context, req model.MovieParams) (*model.MovieResponse, error)
}

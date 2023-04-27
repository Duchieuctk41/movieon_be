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
	GetListContinue(ctx context.Context, req model.MovieParams) (*model.MovieResponse, error)
	GetListMoviesByIdOld(ctx context.Context, listIdOld []string) (*model.MovieResponse, error)

	// user
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	GetOneUser(ctx context.Context, id string) (*model.User, error)
	GetListUser(ctx context.Context, req model.UserParams) (*model.UserResponse, error)
	GetOneUserByEmail(ctx context.Context, email string) (*model.User, error)

	// rating
	CreateRating(ctx context.Context, user *model.Rating) error
	UpdateRating(ctx context.Context, user *model.Rating) error
	DeleteRating(ctx context.Context, id string) error
	GetOneRating(ctx context.Context, id string) (*model.Rating, error)
	GetOneRatingByUM(ctx context.Context, userId string, movieId string) (*model.Rating, error)
	GetListRating(ctx context.Context, req model.RatingParams) (*model.RatingResponse, error)

	// view movie
	CreateViewMovie(ctx context.Context, ob *model.ViewMovie) error
	GetOneViewMovie(ctx context.Context, movieId string, userId string) (*model.ViewMovie, error)
}

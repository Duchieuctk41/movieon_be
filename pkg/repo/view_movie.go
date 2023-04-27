package repo

import (
	"context"
	"movieon_be/pkg/model"
)

func (r *RepoPG) CreateViewMovie(ctx context.Context, ob *model.ViewMovie) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Create(ob).Error
}

func (r *RepoPG) GetOneViewMovie(ctx context.Context, movieId string, userId string) (*model.ViewMovie, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.ViewMovie{}
	if err := tx.Where("movie_id = ? and user_id = ?", movieId, userId).First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

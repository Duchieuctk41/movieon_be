package repo

import (
	"context"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
	"strings"
)

func (r *RepoPG) CreateMovie(ctx context.Context, ob *model.Movie) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Create(ob).Error
}

func (r *RepoPG) UpdateMovie(ctx context.Context, ob *model.Movie) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Where("id = ?", ob.ID).Updates(&ob).Error
}

func (r *RepoPG) DeleteMovie(ctx context.Context, id string) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Where("id = ?", id).Delete(&model.Movie{}).Error
}

func (r *RepoPG) GetOneMovie(ctx context.Context, id string) (*model.Movie, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.Movie{}
	if err := tx.Where("id = ?", id).Find(&rs).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

func (r *RepoPG) GetListMovie(ctx context.Context, req model.MovieParams) (*model.MovieResponse, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.MovieResponse{}
	var err error
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)
	total := new(struct {
		Count int `json:"count"`
	})

	tx = tx.Select("movie.*")

	if req.Day != "" || req.MovieTheaterId != "" {
		tx = tx.Joins("JOIN show ON show.movie_id = movie.id")
		if req.Day != "" {
			tx = tx.Where("show.day = ?", req.Day)
		}

		if req.MovieTheaterId != "" {
			tx = tx.Where("show.movie_theater_id = ?", req.MovieTheaterId)
		}
	}

	if req.Search != "" {
		tx = tx.Where("movie.name like ?", "%"+req.Search+"%")
	}

	if req.Filter != "" {
		filter := strings.Split(req.Filter, ",")
		for i := 0; i < len(filter); i += 2 {
			if i+1 < len(filter) {
				tx = tx.Where(filter[i]+" = ?", filter[i+1])
			}
		}
	}

	switch req.Sort {
	case utils.SORT_CREATED_AT_OLDEST:
		tx = tx.Order("created_at")
	default:
		tx = tx.Order("created_at desc")
	}

	if req.Day != "" || req.MovieTheaterId != "" {
		tx = tx.Group("movie.id")
	}

	if err := tx.Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	//2022-12-18T00:00:00+07:00
	//postman: 2022-12-18T00:00:00 07:00
	if rs.Meta, err = r.GetPaginationInfo("", tx, total.Count, page, pageSize); err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

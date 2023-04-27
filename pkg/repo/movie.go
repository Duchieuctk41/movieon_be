package repo

import (
	"context"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
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

	tx = tx.Select("mo_movie.*").
		Joins("left join view_movie on view_movie.movie_id = mo_movie.id").
		Where("view_movie.movie_id is null")

	switch req.Sort {
	case utils.SORT_CREATED_AT_OLDEST:
		tx = tx.Order("mo_movie.view_count desc, mo_movie.created_at")
	default:
		tx = tx.Order("mo_movie.view_count desc, mo_movie.created_at desc")
	}

	if err := tx.Limit(pageSize).Offset(r.GetOffset(page, pageSize)).Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	if rs.Meta, err = r.GetPaginationInfo("", tx, total.Count, page, pageSize); err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

func (r *RepoPG) GetListContinue(ctx context.Context, req model.MovieParams) (*model.MovieResponse, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.MovieResponse{}
	var err error
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)
	total := new(struct {
		Count int `json:"count"`
	})

	tx = tx.Select("mo_movie.*").
		Joins("inner join view_movie on view_movie.movie_id = mo_movie.id").
		Where("view_movie.user_id = ?", req.UserId)

	if err := tx.Limit(pageSize).Offset(r.GetOffset(page, pageSize)).Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	if rs.Meta, err = r.GetPaginationInfo("", tx, total.Count, page, pageSize); err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

func (r *RepoPG) GetListMoviesByIdOld(ctx context.Context, listIdOld []string) (*model.MovieResponse, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.MovieResponse{}

	tx = tx.Where("id_old in (?)", listIdOld)

	if err := tx.Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

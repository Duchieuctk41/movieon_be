package repo

import (
	"context"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
)

func (r *RepoPG) CreateRating(ctx context.Context, ob *model.Rating) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Create(ob).Error
}

func (r *RepoPG) UpdateRating(ctx context.Context, ob *model.Rating) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Where("id = ?", ob.ID).Updates(&ob).Error
}

func (r *RepoPG) DeleteRating(ctx context.Context, id string) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Where("id = ?", id).Delete(&model.Rating{}).Error
}

func (r *RepoPG) GetOneRating(ctx context.Context, id string) (*model.Rating, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.Rating{}
	if err := tx.Where("id = ?", id).First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r *RepoPG) GetOneRatingByUM(ctx context.Context, userId string, movieId string) (*model.Rating, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.Rating{}
	if err := tx.Where("user_id = ? and movie_id = ?", userId, movieId).First(&rs).Error; err != nil {
		return nil, err
	}

	return &rs, nil
}

func (r *RepoPG) GetListRating(ctx context.Context, req model.RatingParams) (*model.RatingResponse, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.RatingResponse{}
	var err error
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)
	total := new(struct {
		Count int `json:"count"`
	})

	switch req.Sort {
	case utils.SORT_CREATED_AT_OLDEST:
		tx = tx.Order("created_at")
	default:
		tx = tx.Order("created_at desc")
	}

	if err := tx.Limit(pageSize).Offset(r.GetOffset(page, pageSize)).Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	if rs.Meta, err = r.GetPaginationInfo("", tx, total.Count, page, pageSize); err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

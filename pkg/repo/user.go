package repo

import (
	"context"
	"github.com/praslar/cloud0/logger"
	"movieon_be/pkg/model"
	"movieon_be/pkg/utils"
	"strings"
)

func (r *RepoPG) CreateUser(ctx context.Context, ob *model.User) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Create(ob).Error
}

func (r *RepoPG) UpdateUser(ctx context.Context, ob *model.User) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.WithContext(ctx).Where("id = ?", ob.ID).Updates(&ob).Error; err != nil {
		log.WithError(err)
		return err
	}
	return nil
}

func (r *RepoPG) DeleteUser(ctx context.Context, id string) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	return tx.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *RepoPG) GetOneUser(ctx context.Context, id string) (*model.User, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.User{}
	if err := tx.Where("id = ?", id).Take(&rs).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

func (r *RepoPG) GetListUser(ctx context.Context, req model.UserParams) (*model.UserResponse, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.UserResponse{}
	var err error
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)
	total := new(struct {
		Count int `json:"count"`
	})

	if req.Search != "" {
		tx = tx.Where("unaccent(name) ilike %?%", req.Search)
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
	if err := tx.Find(&rs.Data).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	if rs.Meta, err = r.GetPaginationInfo("", tx, total.Count, page, pageSize); err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

func (r *RepoPG) GetOneUserByEmail(ctx context.Context, email string) (*model.User, error) {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	rs := model.User{}
	if err := tx.Where("email = ?", email).Take(&rs).Error; err != nil {
		return nil, r.ReturnErrorInGetFunc(ctx, err, utils.GetCurrentCaller(r, 0))
	}

	return &rs, nil
}

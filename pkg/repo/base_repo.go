package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/cloud0/logger"
	"gorm.io/gorm"
	"math"
	"movieon_be/pkg/utils"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	generalQueryTimeout = 60 * time.Second
	defaultPageSize     = 30
	maxPageSize         = 1000
)

type RepoPG struct {
	DB    *gorm.DB
	debug bool
}

func (r *RepoPG) GetRepo() *gorm.DB {
	return r.DB
}

func NewPGRepo(db *gorm.DB) PGInterface {
	return &RepoPG{DB: db}
}

func (r *RepoPG) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}

func (r *RepoPG) GetPage(page int) int {
	if page == 0 {
		return 1
	}
	return page
}

func (r *RepoPG) GetOffset(page int, pageSize int) int {
	return (page - 1) * pageSize
}

func (r *RepoPG) GetPageSize(pageSize int) int {
	if pageSize == 0 {
		return defaultPageSize
	}
	if pageSize > maxPageSize {
		return maxPageSize
	}
	return pageSize
}

func (r *RepoPG) GetTotalPages(totalRows, pageSize int) int {
	return int(math.Ceil(float64(totalRows) / float64(pageSize)))
}

func (r *RepoPG) GetOrder(sort string) string {
	if sort == "" {
		sort = "created_at desc"
	}
	return sort
}

func (r *RepoPG) GetPaginationInfo(query string, tx *gorm.DB, totalRow, page, pageSize int) (rs ginext.BodyMeta, err error) {
	tm := struct {
		Count int `json:"count"`
	}{}
	if query != "" {
		if err = tx.Raw(query).Scan(&tm).Error; err != nil {
			return nil, err
		}
		totalRow = tm.Count
	}

	return ginext.BodyMeta{
		"page":        page,
		"page_size":   pageSize,
		"total_pages": r.GetTotalPages(totalRow, pageSize),
		"total_rows":  totalRow,
	}, nil
}

func (r *RepoPG) ReturnErrorInGetFunc(ctx context.Context, err error, message string) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.WithError(err).Error(message)
		return ginext.NewError(http.StatusNotFound, err.Error())
	}
	log.WithError(err).Error(message)
	return ginext.NewError(http.StatusInternalServerError, err.Error())
}

func (r *RepoPG) Transaction(ctx context.Context, f func(rp PGInterface) error) (err error) {
	log := logger.WithCtx(ctx, "RepoPG.Transaction")
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()
	// create new instance to run the transaction
	repo := *r
	tx = tx.Begin()
	repo.DB = tx
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = errors.New(fmt.Sprint(r))
			log.WithError(err).Error("error_500: Panic when run Transaction")
			debug.PrintStack()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	err = f(&repo)
	if err != nil {
		log.WithError(err).Error("error_500: Error when run Transaction")
		return err
	}
	return nil
}

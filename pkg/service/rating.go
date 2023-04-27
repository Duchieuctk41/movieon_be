package service

import (
	"context"
	"encoding/csv"
	"github.com/praslar/lib/common"
	"gorm.io/gorm"
	"io"
	"log"
	"movieon_be/pkg/model"
	"movieon_be/pkg/repo"
	"movieon_be/pkg/utils"
	"movieon_be/pkg/valid"
	"os"
	"sync"
	"time"
)

type RatingService struct {
	repo repo.PGInterface
}

type RatingInterface interface {
	Create(ctx context.Context, ob model.RatingRequest) (rs *model.Rating, err error)
	Update(ctx context.Context, ob model.RatingRequest) (rs *model.Rating, err error)
	Delete(ctx context.Context, id string) (err error)
	GetOne(ctx context.Context, id string) (rs *model.Rating, err error)
	GetList(ctx context.Context, req model.RatingParams) (rs *model.RatingResponse, err error)
	CreateOrUpdate(ctx context.Context, ob model.RatingRequest) (rs *model.Rating, err error)
	MashUpload(ctx context.Context) (err error)
}

func NewRatingService(repo repo.PGInterface) RatingInterface {
	return &RatingService{repo: repo}
}

func (s *RatingService) Create(ctx context.Context, req model.RatingRequest) (rs *model.Rating, err error) {

	ob := &model.Rating{}
	common.Sync(req, ob)

	if err := s.repo.CreateRating(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *RatingService) Update(ctx context.Context, req model.RatingRequest) (rs *model.Rating, err error) {
	ob, err := s.repo.GetOneRating(ctx, valid.String(req.ID))
	if err != nil {
		return nil, err
	}

	common.Sync(req, ob)

	if err := s.repo.UpdateRating(ctx, ob); err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *RatingService) Delete(ctx context.Context, id string) (err error) {
	return s.repo.DeleteRating(ctx, id)
}

func (s *RatingService) GetOne(ctx context.Context, id string) (rs *model.Rating, err error) {

	ob, err := s.repo.GetOneRating(ctx, id)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *RatingService) GetList(ctx context.Context, req model.RatingParams) (rs *model.RatingResponse, err error) {

	ob, err := s.repo.GetListRating(ctx, req)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *RatingService) CreateOrUpdate(ctx context.Context, req model.RatingRequest) (rs *model.Rating, err error) {

	ob := &model.Rating{}
	common.Sync(req, ob)

	// check if exist
	rating, err := s.repo.GetOneRatingByUM(ctx, valid.String(req.UserId), valid.String(req.MovieId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := s.repo.CreateRating(ctx, ob); err != nil {
				return nil, err
			}
			return ob, nil
		} else {
			return nil, err
		}
	}

	if rating != nil {
		common.Sync(req, rating)
		if err := s.repo.UpdateRating(ctx, rating); err != nil {
			return nil, err
		}
		return rating, nil
	}

	return ob, nil
}

func (s *RatingService) MashUpload(ctx context.Context) (err error) {
	timeStart := time.Now()

	wg := sync.WaitGroup{}
	ch := make(chan model.RatingCsv)
	wg.Add(2)
	go GetRatingCsv(&wg, ch)
	go s.CreateRatingCsv(&wg, ch)
	wg.Wait()

	timeEnd := time.Now()
	log.Println("done create Ratingcsv: ", timeEnd.Sub(timeStart))

	return nil
}
func GetRatingCsv(wg *sync.WaitGroup, ch chan model.RatingCsv) {
	defer wg.Done()
	defer close(ch)
	file, err := os.Open("assets/csv/ratings.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		item := model.RatingCsv{
			UserId:    record[0],
			MovieId:   record[1],
			Rating:    utils.ConvertStringToFloat64(record[2]),
			Timestamp: record[3],
		}
		ch <- item
	}
}

func (s *RatingService) CreateRatingCsv(wg *sync.WaitGroup, ch chan model.RatingCsv) {
	defer wg.Done()
	for {
		item, ok := <-ch
		if !ok {
			return
		}

		s.repo.CreateRating(context.Background(), &model.Rating{
			UserId:     item.UserId,
			MovieIdOld: item.MovieId,
			Rating:     item.Rating,
			BaseModel: model.BaseModel{
				CreatedAt: utils.ConvertTimeStampToTime(item.Timestamp),
			},
		})
	}
}

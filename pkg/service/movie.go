package service

import (
	"context"
	"encoding/csv"
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/lib/common"
	"gorm.io/gorm"
	"io"
	"log"
	"movieon_be/api"
	"movieon_be/pkg/model"
	"movieon_be/pkg/repo"
	"movieon_be/pkg/utils"
	"movieon_be/pkg/valid"
	"net/http"
	"os"
	"sync"
	"time"
)

type MovieService struct {
	repo      repo.PGInterface
	apiClient *api.ApiClient
}

type MovieInterface interface {
	Create(ctx context.Context, ob model.MovieRequest) (rs *model.Movie, err error)
	Update(ctx context.Context, ob model.MovieRequest) (rs *model.Movie, err error)
	Delete(ctx context.Context, id string) (err error)
	GetOne(ctx context.Context, id string) (rs *model.Movie, err error)
	GetList(ctx context.Context, req model.MovieParams) (rs *model.MovieResponse, err error)
	GetListContinue(ctx context.Context, req model.MovieParams) (rs *model.MovieResponse, err error)
	GetListSuggest(ctx context.Context, idOld string) (rs *model.MovieResponse, err error)
	UpdateViewCount(ctx context.Context, req model.ViewMovieRequest) (rs *model.Movie, err error)
	MashUpload(ctx context.Context) (err error)
}

func NewMovieService(repo repo.PGInterface, api *api.ApiClient) MovieInterface {
	return &MovieService{repo: repo, apiClient: api}
}

func (s *MovieService) Create(ctx context.Context, req model.MovieRequest) (rs *model.Movie, err error) {

	ob := &model.Movie{}
	common.Sync(req, ob)

	// check poster exist. create default image when don't have poster
	if ob.Poster == "" {
		ob.Poster = utils.DEFAULT_IMAGE_URL
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

func (s *MovieService) GetListContinue(ctx context.Context, req model.MovieParams) (rs *model.MovieResponse, err error) {

	ob, err := s.repo.GetListContinue(ctx, req)
	if err != nil {
		return nil, err
	}
	return ob, nil
}

func (s *MovieService) GetListSuggest(ctx context.Context, idOld string) (rs *model.MovieResponse, err error) {

	listIdOld, err := s.apiClient.AiClient.GetListSuggest(ctx, idOld)
	if err != nil {
		return nil, err
	}

	// get list movie by list idOld
	listMovie, err := s.repo.GetListMoviesByIdOld(ctx, listIdOld.Data)
	if err != nil {
		return nil, err
	}

	return listMovie, nil
}

func (s *MovieService) UpdateViewCount(ctx context.Context, req model.ViewMovieRequest) (rs *model.Movie, err error) {

	// get movie by id
	movie, err := s.repo.GetOneMovie(ctx, valid.String(req.MovieId))
	if err != nil {
		return nil, err
	}

	movie.ViewCount = movie.ViewCount + 1
	if err := s.repo.UpdateMovie(ctx, movie); err != nil {
		return nil, err
	}

	// get view_movie by movie_id and user_id
	_, err = s.repo.GetOneViewMovie(ctx, valid.String(req.MovieId), valid.String(req.UserId))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// create view_movie table
			if err := s.repo.CreateViewMovie(ctx, &model.ViewMovie{
				MovieId: valid.String(req.MovieId),
				UserId:  valid.String(req.UserId),
			}); err != nil {
				return nil, err
			}
		} else {
			return nil, ginext.NewError(http.StatusInternalServerError, err.Error())
		}
	}

	return rs, nil
}

func (s *MovieService) MashUpload(ctx context.Context) (err error) {
	timeStart := time.Now()

	wg := sync.WaitGroup{}
	ch := make(chan model.MovieCsv)
	wg.Add(2)
	go GetMovieCsv(&wg, ch)
	go s.CreateMovieCsv(&wg, ch)
	wg.Wait()

	timeEnd := time.Now()
	log.Println("done create moviecsv: ", timeEnd.Sub(timeStart))

	return nil
}

func GetMovieCsv(wg *sync.WaitGroup, ch chan model.MovieCsv) {
	defer wg.Done()
	defer close(ch)
	file, err := os.Open("assets/csv/movies.csv")
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
		item := model.MovieCsv{
			Id:     record[0],
			Title:  record[1],
			Genres: record[2],
		}
		ch <- item
	}
}

func (s *MovieService) CreateMovieCsv(wg *sync.WaitGroup, ch chan model.MovieCsv) {
	defer wg.Done()
	for {
		item, ok := <-ch
		if !ok {
			return
		}

		s.repo.CreateMovie(context.Background(), &model.Movie{
			Name:  item.Title,
			Type:  item.Genres,
			IdOld: item.Id,
		})
	}
}

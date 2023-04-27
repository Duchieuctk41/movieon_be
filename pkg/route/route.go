package route

import (
	"github.com/caarlos0/env/v6"
	"github.com/praslar/cloud0/ginext"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"movieon_be/api"
	"movieon_be/conf"
	"movieon_be/pkg/handlers"
	"movieon_be/pkg/repo"
	service2 "movieon_be/pkg/service"
	"net/http"
)

type extraSetting struct {
	DbDebugEnable bool `env:"DB_DEBUG_ENABLE" envDefault:"true"`
}

type Service struct {
	*conf.BaseApp
	setting *extraSetting
}

// NewService
// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /api/v1
func NewService() *Service {
	s := &Service{
		conf.NewApp("movieon_be", "v1.0"),
		&extraSetting{},
	}
	_ = env.Parse(s.setting)

	db := s.GetDB()
	if s.setting.DbDebugEnable {
		db = db.Debug()
	}
	repoPG := repo.NewPGRepo(db)

	userService := service2.NewUserService(repoPG)
	user := handlers.NewUserHandlers(userService)

	movieService := service2.NewMovieService(repoPG, api.GetApiClient())
	movie := handlers.NewMovieHandlers(movieService)

	ratingService := service2.NewRatingService(repoPG)
	rating := handlers.NewRatingHandlers(ratingService)

	if conf.GetEnv().EnvName == "dev" {
		s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1Api := s.Router.Group("/api/v1")

	// Movie
	v1Api.POST("/admin/movie/create", ginext.WrapHandler(movie.Create))
	v1Api.PUT("/admin/movie/update/:id", ginext.WrapHandler(movie.Update))
	v1Api.DELETE("/admin/movie/delete/:id", ginext.WrapHandler(movie.Delete))
	v1Api.GET("/movie/get-one/:id", ginext.WrapHandler(movie.GetOne))
	v1Api.GET("/movie/get-list", ginext.WrapHandler(movie.GetList))
	v1Api.GET("/movie/get-list-suggest/:id", ginext.WrapHandler(movie.GetListSuggest))
	v1Api.GET("/movie/get-list-con", ginext.WrapHandler(movie.GetListContinue))

	// View
	v1Api.POST("movie/update-view", ginext.WrapHandler(movie.UpdateViewCount))

	// Rating
	v1Api.POST("/rating/create", ginext.WrapHandler(rating.Create))
	v1Api.PUT("/rating/update/:id", ginext.WrapHandler(rating.Update))
	v1Api.DELETE("/rating/delete/:id", ginext.WrapHandler(rating.Delete))
	v1Api.GET("/rating/get-one/:id", ginext.WrapHandler(rating.GetOne))
	v1Api.GET("/rating/get-list", ginext.WrapHandler(rating.GetList))
	v1Api.POST("/rating/cou-rating", ginext.WrapHandler(rating.CreateOrUpdate))

	// MashUpload
	//v1Api.GET("/system/mash-upload/movie-csv", ginext.WrapHandler(movie.MashUpload))
	//v1Api.GET("/system/mash-upload/rating-csv", ginext.WrapHandler(rating.MashUpload))
	//v1Api.GET("/system/mash-upload/user-csv", ginext.WrapHandler(user.MashUpload))

	// Auth
	v1Api.POST("/auth/login", ginext.WrapHandler(user.Login))
	v1Api.POST("/auth/register", ginext.WrapHandler(user.Register))

	// Migrate
	migrateHandler := handlers.NewMigrationHandler(db)
	v1Api.POST("/internal/migrate", migrateHandler.Migrate)

	return s
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

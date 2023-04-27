package model

type ViewMovie struct {
	MovieId string `json:"movie_id" gorm:"movie_id"`
	UserId  string `json:"user_id" gorm:"user_id"`
}

type ViewMovieRequest struct {
	MovieId *string `json:"movie_id"`
	UserId  *string `json:"user_id"`
}

func (ViewMovie) TableName() string {
	return "view_movie"
}

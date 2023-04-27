package model

type Rating struct {
	BaseModel
	UserId     string  `gorm:"column:user_id" json:"user_id"`
	MovieId    string  `gorm:"column:movie_id" json:"movie_id"`
	Rating     float64 `gorm:"type:double;column:rating" json:"rating"`
	UserIdOld  string  `gorm:"column:user_id_old" json:"user_id_old"`
	MovieIdOld string  `gorm:"column:movie_id_old" json:"movie_id_old"`
}

func (Rating) TableName() string {
	return "rating"
}

type RatingRequest struct {
	ID         *string  `json:"id"`
	UserId     *string  `json:"user_id"`
	MovieId    *string  `json:"movie_id"`
	Rating     *float64 `json:"rating"`
	MovieIdOld *string  `json:"movie_id_old"`
}

type RatingParams struct {
	BaseParam
}

type RatingResponse struct {
	Data []Rating               `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

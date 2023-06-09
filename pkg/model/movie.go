package model

type Movie struct {
	BaseModel
	Type        string  `gorm:"type:varchar(250);column:type" json:"type"` // the loai phim
	Cast        string  `gorm:"type:varchar(250);cast" json:"cast"`
	Name        string  `gorm:"type:varchar(250); name;index" json:"name"`
	Spoil       string  `gorm:"type:varchar(500);spoil" json:"spoil"`
	Duration    float64 `gorm:"duration;type:double" json:"duration"`
	ReleaseDate string  `gorm:"type:varchar(100);release_date" json:"release_date"`
	Country     string  `gorm:"type:varchar(100);country" json:"country"`
	Language    string  `gorm:"type:varchar(250);language" json:"language"`
	Rated       string  `gorm:"type:varchar(250);rated" json:"rated"`
	Director    string  `gorm:"type:varchar(250);director" json:"director"`
	Status      string  `gorm:"type:varchar(100);status" json:"status"`
	Poster      string  `gorm:"poster" json:"poster"`
	Trailer     string  `gorm:"trailer" json:"trailer"`
	FormatMovie string  `gorm:"format_movie" json:"format"`
	IdOld       string  `gorm:"id_old" json:"id_old"` // id cua phim tren he thong cu
	ViewCount   int64   `gorm:"view_count" json:"view_count"`
}

func (Movie) TableName() string {
	return "mo_movie"
}

type MovieRequest struct {
	ID          *string  `json:"id"`
	Type        *string  `json:"type"`
	Cast        *string  `json:"cast"`
	Name        *string  `json:"name"`
	Spoil       *string  `json:"spoil"`
	Duration    *float64 `json:"duration"`
	ReleaseDate *string  `json:"release_date"`
	Country     *string  `json:"country"`
	Language    *string  `json:"language"`
	Rated       *string  `json:"rated"`
	Director    *string  `json:"director"`
	Status      *string  `json:"status"`
	Poster      *string  `json:"poster"`
	Trailer     *string  `json:"trailer"`
	FormatMovie *string  `json:"format_movie"`
	IdOld       *string  `json:"id_old"`
	ViewCount   *int64   `json:"view_count"`
}

type MovieParams struct {
	BaseParam
	Day    string `json:"day" form:"day"`
	UserId string `json:"user_id" form:"user_id"`
}

type MovieResponse struct {
	Data []Movie                `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

type MovieSuggestResponse struct {
	IdOld  string `json:"id_old"`
	Name   string `json:"name"`
	Genres string `json:"genres"`
}

type ListSuggestResponse struct {
	Data []string `json:"data"`
}

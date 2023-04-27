package model

type MovieCsv struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Genres string `json:"genres"`
}

type RatingCsv struct {
	UserId    string  `json:"user_id"`
	MovieId   string  `json:"movie_id"`
	Rating    float64 `json:"rating"`
	Timestamp string  `json:"timestamp"`
}

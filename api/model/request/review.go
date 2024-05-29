package request

type NewMovieReview struct {
	MovieID string `json:"movieID" validate:"required"`
	Comment string `json:"comment" validate:"required|max_len:1000"`
	Rating  int    `json:"rating" validate:"required|min:1|max:5"`
}

type UpdateMovieReview struct {
	ID      string  `json:"id" validate:"required"`
	Comment *string `json:"comment,omitempty" validate:"max_len:1000|required_without:Description"`
	Rating  *int    `json:"rating,omitempty" validate:"min:1|max:5"`
}

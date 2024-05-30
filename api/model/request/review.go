package request

type NewMovieReview struct {
	MovieID string `json:"movieID" validate:"required|isNotJustWhitespace"`
	Comment string `json:"comment" validate:"required|max_len:1000|isNotJustWhitespace"`
	Rating  int    `json:"rating" validate:"required|min:1|max:5"`
}

type UpdateMovieReview struct {
	ID      string  `json:"id" validate:"required|isNotJustWhitespace"`
	Comment *string `json:"comment,omitempty" validate:"max_len:1000|required_without:Description|isNotJustWhitespace"`
	Rating  *int    `json:"rating,omitempty" validate:"min:1|max:5"`
}

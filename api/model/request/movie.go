package request

type NewMovie struct {
	Title       string `json:"title" validate:"required|max_len:255"`
	Description string `json:"description" validate:"required|max_len:1000"`
}

type UpdateMovie struct {
	ID          string  `json:"id" validate:"required"`
	Title       *string `json:"title" validate:"max_len:255|required_without:Description"`
	Description *string `json:"description" validate:"max_len:1000"`
}


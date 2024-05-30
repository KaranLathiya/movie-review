package request

type NewMovie struct {
	Title       string `json:"title" validate:"required|max_len:255|isNotJustWhitespace"`
	Description string `json:"description" validate:"required|max_len:1000|isNotJustWhitespace"`
}

type UpdateMovie struct {
	ID          string  `json:"id" validate:"required|isNotJustWhitespace"`
	Title       *string `json:"title" validate:"max_len:255|required_without:Description|isNotJustWhitespace"`
	Description *string `json:"description" validate:"max_len:1000|isNotJustWhitespace"`
}


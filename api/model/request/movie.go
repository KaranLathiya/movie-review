package request

type NewMovie struct {
	Title       string `json:"title" validate:"required|max_len:255"`
	Description string `json:"description" validate:"required|max_len:1000"`
}

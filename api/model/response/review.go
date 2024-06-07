package response

type MovieReviewLimit struct {
	ID     string `json:"id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

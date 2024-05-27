// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Movie struct {
	ID            string    `json:"id"`
	Title         *string   `json:"title,omitempty"`
	Description   *string   `json:"description,omitempty"`
	DirectorID    *string   `json:"directorID,omitempty"`
	CreatedAt     *string   `json:"createdAt,omitempty"`
	UpdatedAt     *string   `json:"updatedAt,omitempty"`
	UpdatedBy     *string   `json:"updatedBy,omitempty"`
	Reviews       []*Review `json:"reviews,omitempty"`
	AverageRating *float64  `json:"averageRating,omitempty"`
}

type MovieReviewNotification struct {
	MovieID    string `json:"movieID"`
	Comment    string `json:"comment"`
	Rating     int    `json:"rating"`
	ReviewerID string `json:"reviewerID"`
	ID         string `json:"id"`
}

type Mutation struct {
}

type Query struct {
}

type Review struct {
	MovieID    *string `json:"movieID,omitempty"`
	Comment    *string `json:"comment,omitempty"`
	Rating     *int    `json:"rating,omitempty"`
	ReviewerID *string `json:"reviewerID,omitempty"`
	ID         *string `json:"id,omitempty"`
}

type Subscription struct {
}

type Token struct {
	AccessToken string `json:"AccessToken"`
}

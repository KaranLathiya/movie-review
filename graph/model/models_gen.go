// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Movie struct {
	ID              string         `json:"id"`
	Title           *string        `json:"title,omitempty"`
	Description     *string        `json:"description,omitempty"`
	DirectorID      *string        `json:"directorID,omitempty"`
	CreatedAt       *string        `json:"createdAt,omitempty"`
	UpdatedAt       *string        `json:"updatedAt,omitempty"`
	UpdatedByUserID *string        `json:"updatedByUserID,omitempty"`
	Reviews         []*MovieReview `json:"reviews,omitempty"`
	AverageRating   *float64       `json:"averageRating,omitempty"`
	Director        *string        `json:"director,omitempty"`
	UpdatedBy       *string        `json:"updatedBy,omitempty"`
}

type MovieReview struct {
	MovieID    *string `json:"movieID,omitempty"`
	Comment    *string `json:"comment,omitempty"`
	Rating     *int    `json:"rating,omitempty"`
	ReviewerID *string `json:"reviewerID,omitempty"`
	ID         *string `json:"id,omitempty"`
	CreatedAt  *string `json:"createdAt,omitempty"`
	UpdatedAt  *string `json:"updatedAt,omitempty"`
	Reviewer   *string `json:"reviewer,omitempty"`
}

type MovieReviewNotification struct {
	ID              string       `json:"id"`
	Title           *string      `json:"title,omitempty"`
	Description     *string      `json:"description,omitempty"`
	DirectorID      *string      `json:"directorID,omitempty"`
	CreatedAt       *string      `json:"createdAt,omitempty"`
	UpdatedAt       *string      `json:"updatedAt,omitempty"`
	UpdatedByUserID *string      `json:"updatedByUserID,omitempty"`
	Review          *MovieReview `json:"review,omitempty"`
	AverageRating   *float64     `json:"averageRating,omitempty"`
	Director        *string      `json:"director,omitempty"`
	UpdatedBy       *string      `json:"updatedBy,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type Subscription struct {
}

type Token struct {
	AccessToken string `json:"AccessToken"`
}

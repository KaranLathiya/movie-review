package repository

import (
	"database/sql"
	"movie-review/api/dal"
	"movie-review/api/model/request"
	"movie-review/constant"
)

type Repositories struct {
	db *sql.DB
}

// InitRepositories should be called in main.go
func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{db: db}
}

type Repository interface {
	UserSignup(user request.UserSignup) error
	UserLogin(user request.UserSignup) (string, error)
	CreateMovie(userId string, movie request.NewMovie) (string, error)
	DeleteMovie(movieID string) error
	UpdateMovie(userID string, movie request.UpdateMovie) error
	CreateMovieReview(userID string, review request.NewMovieReview) (string, error)
	DeleteMovieReview(reviewID string) error
	UpdateMovieReview(userID string, movie request.UpdateMovieReview) error
}

func (r *Repositories) UserSignup(user request.UserSignup) error {
	return dal.UserSignup(r.db, user)
}

func (r *Repositories) UserLogin(user request.UserLogin) (string, error) {
	return dal.UserLogin(r.db, user)
}

func (r *Repositories) CreateMovie(userID string, movie request.NewMovie) (string, error) {
	return dal.CreateMovie(r.db, userID, movie)
}

func (r *Repositories) UpdateMovie(userID string, movie request.UpdateMovie) error {
	return dal.UpdateMovie(r.db, userID, movie)
}

func (r *Repositories) DeleteMovie(movieID string) error {
	return dal.DeleteMovie(r.db, movieID)
}

func (r *Repositories) CreateMovieReview(userID string, review request.NewMovieReview) (string, error) {
	reviewID, err := dal.CreateMovieReview(r.db, userID, review)
	if err != nil {
		return constant.EMPTY_STRING, err
	}
	go dal.UpdateAverageRatingOfMovie(r.db, review.MovieID)
	return reviewID, nil
}

func (r *Repositories) UpdateMovieReview(userID string, movie request.UpdateMovieReview) error {
	movieID, err := dal.UpdateMovieReview(r.db, userID, movie)
	if err != nil {
		return err
	}
	if movie.Rating != nil {
		go dal.UpdateAverageRatingOfMovie(r.db, movieID)
	}
	return nil
}

func (r *Repositories) DeleteMovieReview(userID string, reviewID string) error {
	movieID, err := dal.DeleteMovieReview(r.db, userID, reviewID)
	if err != nil {
		return err
	}
	go dal.UpdateAverageRatingOfMovie(r.db, movieID)
	return nil
}

func (r *Repositories) CheckRoleOfUser(userID string) (string, error) {
	return dal.CheckRoleOfUser(r.db, userID)
}

package repository

import (
	"database/sql"
	"movie-review/api/dal"
	"movie-review/api/model/request"
	"movie-review/constant"
	error_handling "movie-review/error"
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
	tx, err := r.db.Begin()
	if err != nil {
		return constant.EMPTY_STRING, error_handling.InternalServerError
	}
	reviewID, err := dal.CreateMovieReview(tx, userID, review)
	if err != nil {
		defer tx.Rollback()
		return constant.EMPTY_STRING, err
	}
	err = dal.UpdateAverageRatingOfMovie(tx, review.MovieID)
	if err != nil {
		defer tx.Rollback()
		return constant.EMPTY_STRING, err
	}
	err = tx.Commit()
	if err != nil {
		return constant.EMPTY_STRING, error_handling.InternalServerError
	}
	return reviewID, nil
}

func (r *Repositories) UpdateMovieReview(userID string, movie request.UpdateMovieReview) error {
	return dal.UpdateMovieReview(r.db, userID, movie)
}

func (r *Repositories) DeleteMovieReview(userID string, reviewID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	movieID, err := dal.DeleteMovieReview(tx, userID, reviewID)
	if err != nil {
		defer tx.Rollback()
		return err
	}
	err = dal.UpdateAverageRatingOfMovie(tx, movieID)
	if err != nil {
		defer tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) CheckRoleOfUser(userID string) (string, error) {
	return dal.CheckRoleOfUser(r.db, userID)
}

package repository

import (
	"database/sql"
	"movie-review/api/dal"
	"movie-review/api/model/request"
	"movie-review/constant"
	error_handling "movie-review/error"
	"movie-review/graph/model"
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
	CreateMovie(userID string, movie request.NewMovie) (string, error)
	DeleteMovie(movieID string) error
	UpdateMovie(userID string, movie request.UpdateMovie) error
	CreateMovieReview(userID string, review request.NewMovieReview) (string, error)
	DeleteMovieReview(reviewID string) error
	UpdateMovieReview(userID string, movie request.UpdateMovieReview) error
	FetchMovies(movieName string, limit int, offset int) ([]*model.Movie, error)
	FetchMovieReviews(movieID string, limit int, offset int) ([]*model.MovieReview, error)
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
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	err = dal.DeleteMovieReviews(tx, movieID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return error_handling.InternalServerError
		}
		return err
	}
	err = dal.DeleteMovie(tx, movieID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return error_handling.InternalServerError
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) CreateMovieReview(userID string, review request.NewMovieReview) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return constant.EMPTY_STRING, error_handling.InternalServerError
	}
	reviewID, err := dal.CreateMovieReview(tx, userID, review)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return constant.EMPTY_STRING, error_handling.InternalServerError
		}
		return constant.EMPTY_STRING, err
	}
	err = dal.UpdateAverageRatingOfMovie(tx, review.MovieID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return constant.EMPTY_STRING, error_handling.InternalServerError
		}
		return constant.EMPTY_STRING, err
	}
	err = tx.Commit()
	if err != nil {
		return constant.EMPTY_STRING, error_handling.InternalServerError
	}
	return reviewID, nil
}

func (r *Repositories) UpdateMovieReview(userID string, movieReview request.UpdateMovieReview) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	movieID, err := dal.UpdateMovieReview(tx, userID, movieReview)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return error_handling.InternalServerError
		}
		return err
	}
	if movieReview.Rating != nil {
		err = dal.UpdateAverageRatingOfMovie(tx, movieID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return error_handling.InternalServerError
			}
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) DeleteMovieReview(userID string, reviewID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	movieID, err := dal.DeleteMovieReview(tx, userID, reviewID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return error_handling.InternalServerError
		}
		return err
	}
	err = dal.UpdateAverageRatingOfMovie(tx, movieID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return error_handling.InternalServerError
		}
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

func (r *Repositories) FetchMovies(movieName string, limit int, offset int) ([]*model.Movie, error) {
	return dal.FetchMovies(r.db, movieName, limit, offset)
}

func (r *Repositories) FetchMovieReviews(movieID string, limit int, offset int) ([]*model.MovieReview, error) {
	return dal.FetchMovieReviews(r.db, movieID, limit, offset)
}
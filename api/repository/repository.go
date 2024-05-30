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
	SearchMovies(filter *model.MovieSearchFilter, sortBy model.MovieSearchSort, limit int, offset int)
	FetchMovieReviewsByMovieID(movieID string, limit int, offset int) ([]*model.MovieReview, error)
	FetchMovieReviewsUsingDataloader(movieIDs []string, limit string, offset string) ([][]*model.MovieReview, []error)
	FetchMovieByID(movieID string) (*model.Movie, error)
	SearchMovieReviews(filter *model.MovieReviewSearchFilter, sortBy model.MovieReviewSearchSort, limit int, offset int) (*model.MovieReview, error)
	FetchUserDetailsByID(userID string) (*model.UserDetails, error)
}

// for new user signup
func (r *Repositories) UserSignup(user request.UserSignup) error {
	return dal.UserSignup(r.db, user)
}

// for user login using email and password
func (r *Repositories) UserLogin(user request.UserLogin) (string, error) {
	return dal.UserLogin(r.db, user)
}

// to create new movie only by admin
func (r *Repositories) CreateMovie(userID string, movie request.NewMovie) (string, error) {
	return dal.CreateMovie(r.db, userID, movie)
}

// to update movie only by admin
func (r *Repositories) UpdateMovie(userID string, movie request.UpdateMovie) error {
	return dal.UpdateMovie(r.db, userID, movie)
}

// to delete new movie only by admin
func (r *Repositories) DeleteMovie(movieID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	//for deleting movie need to delete all movie reviews first
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

// for creating new movie review
// after creating new movie review update average rating of movie
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
	//after creating new movie review need to update average rating of movie
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

// for updating movie review by review creater only
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
	//no need to update average rating of movie if rating is not updated
	if movieReview.Rating != nil {
		//after updating movie review need to update average rating of movie
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

// for deleteing movie review by review creater only
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
	//after deleting movie review need to update average rating of movie
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

// for fetch role of user
func (r *Repositories) FetchRoleOfUser(userID string) (string, error) {
	return dal.FetchRoleOfUser(r.db, userID)
}

// for seraching movies different filter like title and director and with different sorting options like when created, movie title, movie average rating with ascending and descending options
// for movie title there is full text based search only
func (r *Repositories) SearchMovies(filter *model.MovieSearchFilter, sortBy model.MovieSearchSort, limit int, offset int) ([]*model.Movie, error) {
	return dal.SearchMovies(r.db, filter, sortBy, limit, offset)
}

// for fetching movie by movieID
func (r *Repositories) FetchMovieByID(movieID string) (*model.Movie, error) {
	return dal.FetchMovieByID(r.db, movieID)
}

// for fetching reviews of movie by movieID
func (r *Repositories) FetchMovieReviewsByMovieID(movieID string, limit int, offset int) ([]*model.MovieReview, error) {
	return dal.FetchMovieReviewsByMovieID(r.db, movieID, limit, offset)
}

// for fetching reviews of movies by movieIDs
func (r *Repositories) FetchMovieReviewsUsingDataloader(movieIDs []string) ([][]*model.MovieReview, []error) {
	return dal.FetchMovieReviewsUsingDataloader(r.db, movieIDs)
}

// check the review limit is exceeded or not by the user
// maximum 3 reviews in the last 10 minutes
func (r *Repositories) IsReviewLimitExceeded(userID string) (bool, error) {
	return dal.IsReviewLimitExceeded(r.db, userID)
}

// for searching review by comment(movie review)
// for comment(movie review) there is full text based search only
func (r *Repositories) SearchMovieReviews(filter *model.MovieReviewSearchFilter, sortBy model.MovieReviewSearchSort, limit int, offset int) ([]*model.MovieReview, error) {
	return dal.SearchMovieReviews(r.db, filter, sortBy, limit, offset)
}

// fetching user details by userID
func (r *Repositories) FetchUserDetailsByID(userID string) (*model.UserDetails, error) {
	return dal.FetchUserDetailsByID(r.db, userID)
}

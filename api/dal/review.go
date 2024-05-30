package dal

import (
	"database/sql"
	"fmt"
	"movie-review/api/model/request"
	"movie-review/constant"
	error_handling "movie-review/error"
	"movie-review/graph/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateMovieReview(tx *sql.Tx, userID string, review request.NewMovieReview) (string, error) {
	var reviewID string
	err := tx.QueryRow("INSERT INTO review (movie_id, comment, rating, reviewer_id) VALUES ( $1 , $2 , $3 , $4 ) RETURNING id", review.MovieID, review.Comment, review.Rating, userID).Scan(&reviewID)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return reviewID, nil
}

// this function returns movieID so that using movieID can update average rating of movie
func DeleteMovieReview(tx *sql.Tx, userID string, reviewID string) (string, error) {
	var movieID string
	err := tx.QueryRow("DELETE FROM review WHERE id = $1 AND reviewer_id = $2 RETURNING movie_id", reviewID, userID).Scan(&movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			return constant.EMPTY_STRING, error_handling.MovieReviewDoesNotExist
		}
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return movieID, nil
}

// this function returns mivieID so that using movieID can update average rating of movie
func UpdateMovieReview(tx *sql.Tx, userID string, review request.UpdateMovieReview) (string, error) {
	var movieID string
	var update []string
	var filterArgsList []interface{}

	if review.Comment != nil {
		update = append(update, "comment = ?")
		filterArgsList = append(filterArgsList, review.Comment)
	}
	if review.Rating != nil {
		update = append(update, "rating = ?")
		filterArgsList = append(filterArgsList, review.Rating)
	}
	update = append(update, "updated_at = CURRENT_TIMESTAMP()")
	filterArgsList = append(filterArgsList, review.ID)
	filterArgsList = append(filterArgsList, userID)

	query := fmt.Sprintf("UPDATE review SET %s WHERE id = ? AND reviewer_id = ? RETURNING movie_id", strings.Join(update, " , "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	err := tx.QueryRow(query, filterArgsList...).Scan(&movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			return constant.EMPTY_STRING, error_handling.MovieReviewDoesNotExist
		}
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return movieID, nil
}

// for deleting all movie reviews of particular movie
func DeleteMovieReviews(tx *sql.Tx, movieID string) error {
	_, err := tx.Exec("DELETE FROM review WHERE movie_id = $1", movieID)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	return nil
}

func FetchMovieReviewsByMovieID(db *sql.DB, movieID string, limit int, offset int) ([]*model.MovieReview, error) {
	query := "SELECT r.id, r.movie_id, r.reviewer_id, r.rating, r.comment, r.created_at, r.updated_at, u.full_name AS reviewer_name FROM review r INNER JOIN users u ON r.reviewer_id = u.id WHERE r.movie_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ? "
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, movieID, limit, offset)
	if err != nil {
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	var movieReviews []*model.MovieReview
	for rows.Next() {
		var movieReview model.MovieReview
		err = rows.Scan(&movieReview.ID, &movieReview.MovieID, &movieReview.ReviewerID, &movieReview.Rating, &movieReview.Comment, &movieReview.CreatedAt, &movieReview.UpdatedAt, &movieReview.Reviewer)
		if err != nil {
			return nil, error_handling.InternalServerError
		}
		movieReviews = append(movieReviews, &movieReview)
	}
	if err = rows.Close(); err != nil {
		return nil, error_handling.InternalServerError
	}
	return movieReviews, nil
}

func FetchMovieReviewsUsingDataloader(db *sql.DB, movieIDs []string) ([][]*model.MovieReview, []error) {
	sqlQuery := "SELECT r.id, r.movie_id, r.reviewer_id, r.rating, r.comment, r.created_at, r.updated_at, u.full_name AS reviewer_name FROM review r INNER JOIN users u ON r.reviewer_id = u.id WHERE r.movie_id IN (?) ORDER BY r.created_at DESC"
	sqlQuery, arguments, err := sqlx.In(sqlQuery, movieIDs)
	if err != nil {
		return nil, []error{error_handling.InternalServerError}
	}
	sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
	rows, err := db.Query(sqlQuery, arguments...)
	if err != nil {
		return nil, []error{error_handling.InternalServerError}
	}
	movieReviewMap := map[string][]*model.MovieReview{}
	for rows.Next() {
		movieReview := model.MovieReview{}
		if err = rows.Scan(&movieReview.ID, &movieReview.MovieID, &movieReview.ReviewerID, &movieReview.Rating, &movieReview.Comment, &movieReview.CreatedAt, &movieReview.UpdatedAt, &movieReview.Reviewer); err != nil {
			return nil, []error{error_handling.InternalServerError}
		}
		if _, ok := movieReviewMap[*movieReview.MovieID]; ok {
			movieReviewMap[*movieReview.MovieID] = append(movieReviewMap[*movieReview.MovieID], &movieReview)
		} else {
			movieReviewArr := []*model.MovieReview{}
			movieReviewArr = append(movieReviewArr, &movieReview)
			movieReviewMap[*movieReview.MovieID] = movieReviewArr
		}
	}
	if err = rows.Close(); err != nil {
		return nil, []error{error_handling.InternalServerError}
	}
	movieReviews := make([][]*model.MovieReview, len(movieIDs))
	for i, id := range movieIDs {
		movieReviews[i] = movieReviewMap[id]
		i++
	}
	return movieReviews, nil
}

func SearchMovieReviews(db *sql.DB, filter *model.MovieReviewSearchFilter, sortBy model.MovieReviewSearchSort, limit int, offset int) ([]*model.MovieReview, error) {
	var args []interface{}
	var where, orderBy []string
	var whereKeyword string
	if filter != nil {
		if filter.Reviewer != nil {
			args = append(args, *filter.Reviewer)
			where = append(where, " u.full_name ILIKE '%' || ? || '%' ")
		}
		if filter.Comment != nil {
			args = append(args, *filter.Comment, *filter.Comment)
			//fulltext based search condition only
			where = append(where, " comment_tsvector @@ PLAINTO_TSQUERY(?) ")
			orderBy = append(orderBy, " TS_RANK(comment_tsvector, PLAINTO_TSQUERY(?)) ")
		}
	}

	switch sortBy.String() {
	case model.MovieReviewSearchSortNewest.String():
		orderBy = append(orderBy, " created_at DESC ")
	case model.MovieReviewSearchSortOldest.String():
		orderBy = append(orderBy, " created_at ASC ")
	case model.MovieReviewSearchSortRatingAsc.String():
		orderBy = append(orderBy, " rating ASC ")
	case model.MovieReviewSearchSortRatingDesc.String():
		orderBy = append(orderBy, " rating DESC ")
	}

	if len(where) > 0 {
		whereKeyword = " WHERE "
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf("SELECT r.id, r.movie_id, r.reviewer_id, r.rating, r.comment, r.created_at, r.updated_at, u.full_name AS reviewer_name FROM review r INNER JOIN users u ON r.reviewer_id = u.id  %v %s ORDER BY %s LIMIT ? OFFSET ? ", whereKeyword, strings.Join(where, " AND "), strings.Join(orderBy, " , "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	var movieReviews []*model.MovieReview
	for rows.Next() {
		var movieReview model.MovieReview
		err = rows.Scan(&movieReview.ID, &movieReview.MovieID, &movieReview.ReviewerID, &movieReview.Rating, &movieReview.Comment, &movieReview.CreatedAt, &movieReview.UpdatedAt, &movieReview.Reviewer)
		if err != nil {
			return nil, error_handling.InternalServerError
		}
		movieReviews = append(movieReviews, &movieReview)
	}
	if err = rows.Close(); err != nil {
		return nil, error_handling.InternalServerError
	}
	return movieReviews, nil
}

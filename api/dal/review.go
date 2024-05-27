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
	update = append(update, "updated_at = current_timestamp()")
	filterArgsList = append(filterArgsList, review.ID)
	filterArgsList = append(filterArgsList, userID)

	query := fmt.Sprintf("UPDATE review SET %s WHERE id = ? AND reviewer_id = ? RETURNING movie_id", strings.Join(update, " , "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	err := tx.QueryRow(query, filterArgsList...).Scan(&movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			return constant.EMPTY_STRING, error_handling.InvalidDetails
		}
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return constant.EMPTY_STRING, nil
}

func DeleteMovieReviews(tx *sql.Tx, movieID string) error {
	_, err := tx.Exec("DELETE FROM review WHERE movie_id = $1", movieID)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	return nil
}

func FetchMovieReviews(db *sql.DB, movieID string, limit int, offset int) ([]*model.MovieReview, error) {
	query := "SELECT r.id, r.reviewer_id, r.rating, r.comment, r.created_at, r.updated_at, CONCAT(u.first_name, ' ', u.last_name) AS reviewer_name FROM review r LEFT JOIN users u ON r.reviewer_id = u.id WHERE r.movie_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ? "
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, movieID, limit, offset)
	if err != nil {
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	var movieReviews []*model.MovieReview
	for rows.Next() {
		var movieReview model.MovieReview
		err = rows.Scan(&movieReview.ID, &movieReview.ReviewerID, &movieReview.Rating, &movieReview.Comment, &movieReview.CreatedAt, &movieReview.UpdatedAt, &movieReview.Reviewer)
		if err != nil {
			return nil, error_handling.InternalServerError
		}
		movieReviews = append(movieReviews, &movieReview)
	}
	return movieReviews, nil
}

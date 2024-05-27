package dal

import (
	"database/sql"
	"fmt"
	"movie-review/api/model/request"
	"movie-review/constant"
	error_handling "movie-review/error"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateMovie(db *sql.DB, userID string, movie request.NewMovie) (string, error) {
	var movieID string
	err := db.QueryRow("INSERT INTO movie (title, director_id, description) VALUES ( $1 , $2 , $3 ) RETURNING id", movie.Title, userID, movie.Description).Scan(&movieID)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return movieID, nil
}

func DeleteMovie(tx *sql.Tx, movieID string) error {
	result, err := tx.Exec("DELETE FROM movie WHERE id = $1", movieID)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return error_handling.InternalServerError
	}
	if rowsAffected == 0 {
		return error_handling.MovieDoesNotExist
	}
	return nil
}

func UpdateMovie(db *sql.DB, userID string, movie request.UpdateMovie) error {
	var update []string
	var filterArgsList []interface{}

	if movie.Title != nil {
		update = append(update, "title = ?")
		filterArgsList = append(filterArgsList, movie.Title)
	}
	if movie.Description != nil {
		update = append(update, "description = ?")
		filterArgsList = append(filterArgsList, movie.Description)
	}
	update = append(update, "updated_by = ?")
	filterArgsList = append(filterArgsList, userID)
	update = append(update, "updated_at = current_timestamp()")
	filterArgsList = append(filterArgsList, movie.ID)

	query := fmt.Sprintf("UPDATE movie SET %s WHERE id = ?", strings.Join(update, " , "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	result, err := db.Exec(query, filterArgsList...)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return error_handling.InternalServerError
	}
	if rowsAffected == 0 {
		return error_handling.MovieDoesNotExist
	}
	return nil
}

func UpdateAverageRatingOfMovie(tx *sql.Tx, movieID string) error {
	result, err := tx.Exec(`
	WITH avg_rating AS (
        SELECT
            AVG(rating) AS avg_rating
        FROM
            review
        WHERE
            movie_id = $1
    )
    UPDATE movie
    SET average_rating = (SELECT avg_rating FROM avg_rating)
    WHERE id = $1;`, movieID)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return error_handling.InternalServerError
	}
	if rowsAffected == 0 {
		return error_handling.MovieDoesNotExist
	}
	return nil
}

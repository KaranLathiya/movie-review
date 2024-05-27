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

func FetchMovies(db *sql.DB, movieName string, limit int, offset int) ([]*model.Movie, error) {
	query := "SELECT m.id, m.title, m.description, m.director_id, m.created_at, m.updated_at, m.updated_by, CONCAT(u1.first_name, ' ', u1.last_name) AS director_name, CONCAT(u2.first_name, ' ', u2.last_name) AS updater_name FROM movie m LEFT JOIN users u1 ON m.director_id = u1.id LEFT JOIN users u2 ON m.updated_by = u2.id WHERE m.title ILIKE '%' || ? || '%' ORDER BY created_at DESC LIMIT ? OFFSET ? "
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, movieName, limit, offset)
	if err != nil {
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	var movies []*model.Movie
	for rows.Next() {
		var movie model.Movie
		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.DirectorID, &movie.CreatedAt, &movie.UpdatedAt, &movie.UpdatedByUserID, &movie.Director, &movie.UpdatedBy)
		if err != nil {
			return nil, error_handling.InternalServerError
		}
		movies = append(movies, &movie)
	}
	return movies, nil
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



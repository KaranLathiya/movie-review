package dal

import (
	"database/sql"
	"movie-review/api/constant"
	"movie-review/api/model/request"
	error_handling "movie-review/error"
)

func CreateMovie(db *sql.DB, userID string, movie request.NewMovie) (string, error) {
	var movieID string
	err := db.QueryRow("INSERT INTO movie (title, director_id, description) VALUES ( $1 , $2 , $3 ) RETURNING id", movie.Title, userID, movie.Description).Scan(&movieID)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.DatabaseErrorShow(err)
	}
	return movieID, nil
}

func DeleteMovie(db *sql.DB, movieID string) error {
	_,err := db.Exec("DELETE FROM movie WHERE id = $1", movieID)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
	}
	return nil
}
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
	update = append(update, "updated_at = CURRENT_TIMESTAMP()")
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

func SearchMovies(db *sql.DB, filter *model.MovieSearchFilter, sortBy *model.MovieSearchSort, limit int, offset int) ([]*model.Movie, error) {
	var args []interface{}
	var where, orderBy []string
	var whereKeyword string
	if filter != nil {
		if filter.Director != nil {
			args = append(args, filter.Director)
			where = append(where, " u1.full_name ILIKE '%' || ? || '%' ")
		}
		if filter.Title != nil {
			args = append(args, filter.Title, filter.Title)
			where = append(where, " title_tsvector @@ PLAINTO_TSQUERY(?) ")
			orderBy = append(orderBy, " TS_RANK(title_tsvector, PLAINTO_TSQUERY(?)) ")
		}
	}

	switch sortBy.String() {
	case model.MovieSearchSortNewest.String():
		orderBy = append(orderBy, " created_at DESC ")
	case model.MovieSearchSortOldest.String():
		orderBy = append(orderBy, " created_at ASC ")
	case model.MovieSearchSortTitleDesc.String():
		orderBy = append(orderBy, " LOWER(title) DESC ")
	case model.MovieSearchSortTitleAsc.String():
		orderBy = append(orderBy, " LOWER(titlr) ASC ")
	case model.MovieSearchSortAverageRatingAsc.String():
		orderBy = append(orderBy, " average_rating ASC ")
	case model.MovieSearchSortAverageRatingDesc.String():
		orderBy = append(orderBy, " average_rating DESC ")
	}

	if len(where) > 0 {
		whereKeyword = " WHERE "
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf("SELECT m.id, m.title, m.description, m.director_id, m.created_at, m.updated_at, m.updated_by, m.average_rating, u1.full_name AS director_name, u2.full_name AS updater_name FROM movie m INNER JOIN users u1 ON m.director_id = u1.id LEFT JOIN users u2 ON m.updated_by = u2.id %v %s ORDER BY %s LIMIT ? OFFSET ? ", whereKeyword, strings.Join(where, " AND "), strings.Join(orderBy, " , "))
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	var movies []*model.Movie
	for rows.Next() {
		var movie model.Movie
		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.DirectorID, &movie.CreatedAt, &movie.UpdatedAt, &movie.UpdatedByUserID, &movie.AverageRating, &movie.Director, &movie.UpdatedBy)
		if err != nil {
			return nil, error_handling.InternalServerError
		}
		movies = append(movies, &movie)
	}
	if err = rows.Close(); err != nil {
		return nil, error_handling.InternalServerError
	}
	return movies, nil
}

func FetchMovieByID(db *sql.DB, movieID string) (*model.Movie, error) {
	var movie model.Movie
	err := db.QueryRow("SELECT m.id, m.title, m.description, m.director_id, m.created_at, m.updated_at, m.updated_by, m.average_rating, u1.full_name AS director_name, u2.full_name AS updater_name FROM movie m INNER JOIN users u1 ON m.director_id = u1.id LEFT JOIN users u2 ON m.updated_by = u2.id WHERE m.id = $1 ;", movieID).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.DirectorID, &movie.CreatedAt, &movie.UpdatedAt, &movie.UpdatedByUserID, &movie.AverageRating, &movie.Director, &movie.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, error_handling.MovieDoesNotExist
		}
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	return &movie, nil
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

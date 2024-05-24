package dal

import (
	"database/sql"
	"movie-review/api/model/request"
	"movie-review/utils"

	error_handling "movie-review/error"

	"github.com/lib/pq"
)

func UserSignup(db *sql.DB, user request.UserSignup) error {
	_, err := db.Exec("INSERT INTO users (email,password,first_name,last_name) VALUES ( $1 , $2 , $3 , $4 )", user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code == "23505" {
				return error_handling.UserAlreadyExist
			}
			return err
		}
	}
	return nil
}

func UserLogin(db *sql.DB, user request.UserLogin) (string, error) {
	var id, password string
	err := db.QueryRow("SELECT id, password from users WHERE email = $1", user.Email).Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", error_handling.UserDoesNotExist
		}
		return "", err
	}
	passwordMatch := utils.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if passwordMatch {
		return id, nil
	} else {
		return "", error_handling.InvalidDetails
	}
}

package dal

import (
	"database/sql"
	"movie-review/api/model/request"
	"movie-review/constant"
	"movie-review/utils"

	error_handling "movie-review/error"
)

func UserSignup(db *sql.DB, user request.UserSignup) error {
	_, err := db.Exec("INSERT INTO users (email, password ,first_name ,last_name) VALUES ( $1 , $2 , $3 , $4 )", user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return error_handling.DatabaseErrorHandling(err)
	}
	return nil
}

func UserLogin(db *sql.DB, user request.UserLogin) (string, error) {
	var id, password string
	err := db.QueryRow("SELECT id, password from users WHERE email = $1", user.Email).Scan(&id, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return constant.EMPTY_STRING, error_handling.UserDoesNotExist
		}
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	passwordMatch := utils.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if !passwordMatch {
		return constant.EMPTY_STRING, error_handling.InvalidDetails
	}
	return id, nil
}

func CheckRoleOfUser(db *sql.DB, userID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role from users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return constant.EMPTY_STRING, error_handling.UserDoesNotExist
		}
		return constant.EMPTY_STRING, error_handling.DatabaseErrorHandling(err)
	}
	return role, nil
}

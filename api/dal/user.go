package dal

import (
	"database/sql"
	"movie-review/constant"
	"movie-review/api/model/request"
	"movie-review/utils"

	error_handling "movie-review/error"
)

func UserSignup(db *sql.DB, user request.UserSignup) error {
	_, err := db.Exec("INSERT INTO users (email, password ,first_name ,last_name) VALUES ( $1 , $2 , $3 , $4 )", user.Email, user.Password, user.FirstName, user.LastName)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
	}
	return nil
}

func UserLogin(db *sql.DB, user request.UserLogin) (string, error) {
	var id, password string
	err := db.QueryRow("SELECT id, password from users WHERE email = $1", user.Email).Scan(&id, &password)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.DatabaseErrorShow(err)
	}
	passwordMatch := utils.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if passwordMatch {
		return id, nil
	} else {
		return constant.EMPTY_STRING, error_handling.InvalidDetails
	}
}

func CheckRoleOfUser(db *sql.DB, userID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role from users WHERE id = $1", userID).Scan(&role)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.DatabaseErrorShow(err)
	}
	return role, nil
}

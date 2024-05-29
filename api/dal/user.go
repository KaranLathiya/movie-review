package dal

import (
	"database/sql"
	"movie-review/api/model/request"
	"movie-review/constant"
	"movie-review/graph/model"
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

func FetchUserDetailsByID(db *sql.DB, userID string) (*model.UserDetails, error) {
	var userDetails model.UserDetails
	err := db.QueryRow("SELECT email, first_name, last_name from users WHERE id = $1", userID).Scan(&userDetails.Email,&userDetails.FirstName,&userDetails.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, error_handling.UserDoesNotExist
		}
		return nil, error_handling.DatabaseErrorHandling(err)
	}
	return &userDetails, nil
}

func IsReviewLimitExceeded(db *sql.DB, userID string) (bool, error) {
	var numberOfReviews int
	err := db.QueryRow("SELECT COUNT(*) FROM review WHERE reviewer_id = $1 AND created_at >= CURRENT_TIMESTAMP() - INTERVAL '10 minutes';", userID).Scan(&numberOfReviews)
	if err != nil {
		return false, error_handling.DatabaseErrorHandling(err)
	}
	if numberOfReviews > 3 {
		return true, nil
	}
	return false, nil
}

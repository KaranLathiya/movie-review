package utils

import (
	"movie-review/config"
	"movie-review/constant"
	error_handling "movie-review/error"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
}

func CreateJWT(tokenExpiryTime time.Time, userID string) (string, error) {
	jwtKey := []byte(config.ConfigVal.JWTKey)
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return constant.EMPTY_STRING, error_handling.JWTTokenGenerateError
	}
	return tokenString, nil
}

func VerifyJWT(token string) (string, error) {
	jwtKey := []byte(config.ConfigVal.JWTKey)
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return constant.EMPTY_STRING, error_handling.JWTErrSignatureInvalid
		}
		return constant.EMPTY_STRING, error_handling.CustomError{StatusCode: http.StatusInternalServerError, ErrorMessage: err.Error()}
	}

	if !tkn.Valid {
		return constant.EMPTY_STRING, error_handling.JWTTokenInvalid
	}
	return claims.Subject, nil
}

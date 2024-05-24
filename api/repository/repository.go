package repository

import (
	"database/sql"
	"movie-review/api/dal"
	"movie-review/api/model/request"
)

type Repositories struct {
	db *sql.DB
}

// InitRepositories should be called in main.go
func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{db: db}
}

type Repository interface {
	UserSignup(user request.UserSignup) error
	UserLogin(user request.UserSignup) (string, error)
	CreateMovie(userId string, movie request.NewMovie) (string, error)
	DeleteMovie(movieID string) error
}

func (r *Repositories) UserSignup(user request.UserSignup) error {
	return dal.UserSignup(r.db, user)
}

func (r *Repositories) UserLogin(user request.UserLogin) (string, error) {
	return dal.UserLogin(r.db, user)
}

func (r *Repositories) CreateMovie(userID string, movie request.NewMovie) (string, error) {
	return dal.CreateMovie(r.db, userID, movie)
}

func (r *Repositories) DeleteMovie(movieID string) error{
	return dal.DeleteMovie(r.db, movieID)
}

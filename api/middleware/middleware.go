package middleware

import (
	"context"
	"movie-review/api/repository"
	"net/http"
)

var RepoCtxKey = &contextKey{nil}

type contextKey struct {
	repos *repository.Repositories
}

func AddRepoToContext(repos *repository.Repositories) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = context.WithValue(r.Context(), RepoCtxKey, repos)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

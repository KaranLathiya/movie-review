package middleware

import (
	"context"
	"movie-review/api/repository"
	"net/http"
)

var RepoCtxKey = &repoContextKey{nil}

type repoContextKey struct {
	repos *repository.Repositories
}

var UserCtxKey = &userContextKey{""}

type userContextKey struct {
	token string
}

func AddRepoToContext(repos *repository.Repositories) (func(http.Handler) http.Handler) {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = context.WithValue(r.Context(), RepoCtxKey, repos)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		
		// Allow unauthenticated users in
		if accessToken == "" {
			ctx := context.WithValue(r.Context(), UserCtxKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, accessToken)

		// and call the next with our new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


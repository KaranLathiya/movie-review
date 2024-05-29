package middleware

import (
	"context"
	"movie-review/api/repository"
	"movie-review/constant"
	"net/http"
)

func AddRepoToContext(repos *repository.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = context.WithValue(r.Context(), constant.RepoCtxKey, repos)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(constant.HEADER_KEY_AUTHORIZATION)

		// Allow unauthenticated users in
		if accessToken == "" {
			ctx := context.WithValue(r.Context(), constant.AccessTokenCtxKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), constant.AccessTokenCtxKey, accessToken)

		// and call the next with our new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

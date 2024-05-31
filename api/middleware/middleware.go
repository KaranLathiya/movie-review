package middleware

import (
	"context"
	"movie-review/api/repository"
	"movie-review/constant"
	"net/http"
)

// add repo in the context
func AddRepoToContext(repos *repository.Repositories) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = context.WithValue(r.Context(), constant.RepoCtxKey, repos)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// add accesstoken in the context (from header)
func AddAccessTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(constant.HEADER_KEY_AUTHORIZATION)

		// Allow unauthenticated users in
		if accessToken == constant.EMPTY_STRING {
			ctx := context.WithValue(r.Context(), constant.AccessTokenCtxKey, constant.EMPTY_STRING)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), constant.AccessTokenCtxKey, accessToken)

		// and call the next with our new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

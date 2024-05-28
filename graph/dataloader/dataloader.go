package dataloader

import (
	"context"
	"movie-review/api/middleware"
	"movie-review/api/repository"
	"net/http"
	"time"
)

type Loaders struct {
	ReviewLoader *ReviewLoader
}

var LoaderCtxKey = &LoaderContextKey{"loader"}

type LoaderContextKey struct {
	loader string
}

func DataloaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repo, _ := ctx.Value(middleware.RepoCtxKey).(*repository.Repositories)
		loader := Loaders{}
		loader.ReviewLoader = &ReviewLoader{
			wait:     1 * time.Millisecond,
			maxBatch: 100,
			fetch:    repo.FetchMovieReviewsUsingDataloader,
		}
		ctx = context.WithValue(ctx, LoaderCtxKey, &loader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

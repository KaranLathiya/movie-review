package dataloader

import (
	"context"
	"movie-review/api/repository"
	"movie-review/constant"
	"net/http"
	"time"
)

type Loaders struct {
	ReviewLoader      *ReviewLoader
	ReviewLimitLoader *ReviewLimitLoader
}

// add loader in the context
func DataloaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		repo := ctx.Value(constant.RepoCtxKey).(*repository.Repositories)
		loader := Loaders{}
		loader.ReviewLoader = &ReviewLoader{
			wait:     1 * time.Millisecond,
			maxBatch: 100,
			fetch:    repo.FetchMovieReviewsUsingDataloader,
		}
		loader.ReviewLimitLoader = &ReviewLimitLoader{
			wait:     1 * time.Millisecond,
			maxBatch: 100,
			fetch:    repo.FetchLimitedMovieReviewsUsingDataloader,
		}
		ctx = context.WithValue(ctx, constant.LoaderCtxKey, &loader)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

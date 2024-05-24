package graph

import (
	"context"
	"database/sql"
	"movie-review/api/middleware"
	"movie-review/utils"

	error_handling "movie-review/error"
	"github.com/99designs/gqlgen/graphql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct{}

func NewRootResolvers(db *sql.DB) Config {
	config := Config{
		Resolvers: &Resolver{},
	}

	// Schema Directive
	config.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		accessToken := ctx.Value(middleware.UserCtxKey).(string)
		if accessToken != "" {
			userID, errorMessage := utils.VerifyJWT(accessToken)
			if err != nil {
				ctx := context.WithValue(ctx, middleware.UserCtxKey, userID)
				return next(ctx)
			} else {
				return nil, errorMessage
			}
		} else {
			return nil, error_handling.HeaderDataMissing
		}
	}
	return config
}

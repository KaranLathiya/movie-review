package graph

import (
	"context"
	"movie-review/constant"
	"movie-review/api/middleware"
	"movie-review/api/repository"
	"movie-review/utils"

	error_handling "movie-review/error"

	"github.com/99designs/gqlgen/graphql"

	_ "github.com/99designs/gqlgen/graphql/introspection"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct{}

func NewRootResolvers(repo *repository.Repositories) Config {
	config := Config{
		Resolvers: &Resolver{},
	}

	// Schema Directive
	config.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		accessToken := ctx.Value(middleware.UserCtxKey).(string)
		if accessToken != "" {
			userID, err := utils.VerifyJWT(accessToken)
			if err != nil {
				return nil, err
			} else {
				ctx := context.WithValue(ctx, middleware.UserCtxKey, userID)
				return next(ctx)
			}
		} else {
			return nil, error_handling.HeaderDataMissing
		}
	}

	config.Directives.IsAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		userID, _ := ctx.Value(middleware.UserCtxKey).(string)
		roleOfUser, err := repo.CheckRoleOfUser(userID)
		if err != nil {
			if err == error_handling.NoRowsError {
				return nil, error_handling.UserDoesNotExist
			}
			return nil, err
		}
		if roleOfUser != constant.ADMIN_ROLE{
			return nil, error_handling.AdminAccessRights
		}
		return next(ctx)
	}
	return config
}

package graph

import (
	"context"
	"movie-review/api/repository"
	"movie-review/constant"
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

	//to verify the user
	//check the accesstoken is valid or not if valid then parse userID otherwies throw error
	config.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		accessToken := ctx.Value(constant.AccessTokenCtxKey).(string)
		if accessToken != constant.EMPTY_STRING {
			userID, err := utils.VerifyJWT(accessToken)
			if err != nil {
				return nil, err
			}
			ctx := context.WithValue(ctx, constant.UserIDCtxKey, userID)
			return next(ctx)
		}
		return nil, error_handling.HeaderDataMissing

	}

	//check user is admin or not
	//for admin authorized actions allow only if admin otherwise throw error
	config.Directives.IsAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		userID, _ := ctx.Value(constant.UserIDCtxKey).(string)
		roleOfUser, err := repo.FetchRoleOfUser(userID)
		if err != nil {
			if err == error_handling.NoRowsError {
				return nil, error_handling.UserDoesNotExist
			}
			return nil, err
		}
		if roleOfUser != constant.ADMIN_ROLE {
			return nil, error_handling.AdminAccessRights
		}
		return next(ctx)
	}
	return config
}

package graph

import (
	"database/sql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate
var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	token string
}

type Resolver struct{}

func NewRootResolvers(db *sql.DB) Config {
	config := Config{
		Resolvers: &Resolver{},
	}

	// Schema Directive
	// c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	// 	authorizationKey := ctx.Value(UserCtxKey).(string)
	// 	if authorizationKey != "" {
	// 		ok, errorMessage := utils.VerifyJWT()
	// 		if ok {
	// 			return next(ctx)
	// 		} else {
	// 			return nil, errors.New(errorMessage)
	// 		}
	// 	} else {
	// 		fmt.Println("no autho")
	// 		return nil, errors.New("no authorization key")
	// 	}
	// }

	return config
}

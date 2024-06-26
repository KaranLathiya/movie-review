package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"movie-review/api/model/request"
	"movie-review/api/repository"
	"movie-review/constant"
	error_handling "movie-review/error"
	"movie-review/graph/model"
	"movie-review/utils"
	"time"
)

// UserSignup is the resolver for the UserSignup field.
func (r *mutationResolver) UserSignup(ctx context.Context, input request.UserSignup) (string, error) {
	err := utils.ValidateStruct(input, nil)
	if err != nil {
		return constant.EMPTY_STRING, err
	}
	hashedPassword, err := utils.Bcrypt(input.Password)
	if err != nil {
		return constant.EMPTY_STRING, err
	}
	input.Password = hashedPassword
	repo := ctx.Value(constant.RepoCtxKey).(*repository.Repositories)
	err = repo.UserSignup(input)
	if err != nil {
		if err == error_handling.UniqueKeyConstraintError {
			return constant.EMPTY_STRING, error_handling.UserAlreadyExist
		}
		return constant.EMPTY_STRING, err
	}
	return constant.SIGNUP_SUCCESS, nil
}

// UserLogin is the resolver for the UserLogin field.
func (r *mutationResolver) UserLogin(ctx context.Context, input request.UserLogin) (*model.Token, error) {
	err := utils.ValidateStruct(input, nil)
	if err != nil {
		return nil, err
	}
	repo := ctx.Value(constant.RepoCtxKey).(*repository.Repositories)
	userID, err := repo.UserLogin(input)
	if err != nil {
		return nil, err
	}
	accessToken, err := utils.CreateJWT(time.Now().UTC().Add(time.Minute*time.Duration(60)), userID)
	if err != nil {
		return nil, err
	}
	return &model.Token{AccessToken: accessToken}, nil
}

// FetchCurrentUserDetails is the resolver for the FetchCurrentUserDetails field.
func (r *queryResolver) FetchCurrentUserDetails(ctx context.Context) (*model.UserDetails, error) {
	userID := ctx.Value(constant.UserIDCtxKey).(string)
	repo := ctx.Value(constant.RepoCtxKey).(*repository.Repositories)
	userDetails, err := repo.FetchUserDetailsByID(userID)
	if err != nil {
		return nil, err
	}
	return userDetails, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

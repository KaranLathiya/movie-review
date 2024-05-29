package constant

type contextKey interface{}

var (
	AccessTokenCtxKey = contextKey("accessToken")
	UserIDCtxKey      = contextKey("userID")
	RepoCtxKey        = contextKey("repository")
	LoaderCtxKey      = contextKey("loader")
)

const (
	EMPTY_STRING = ""

	SIGNUP_SUCCESS       = "User signed up successfully."
	MOVIE_UPDATED        = "Movie updated successfully."
	MOVIE_DELETED        = "Movie deleted successfully."
	MOVIE_REVIEW_UPDATED = "Movie review updated successfully."
	MOVIE_REVIEW_DELETED = "Movie review deleted successfully."

	ADMIN_ROLE = "ADMIN"
	USER_ROLE  = "USER"

	HEADER_KEY_AUTHORIZATION = "Authorization"
)

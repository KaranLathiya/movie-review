package error

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gookit/validate"
	"github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func init() {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})
	//passwordRegex validator to check that field must have one special character, one capital letter, one small character, one digit 
	validate.AddValidator("passwordRegex", func(val any) bool {
		conditions := []string{".{7,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
		for _, condition := range conditions {
			conditionMatch, _ := regexp.MatchString(condition, val.(string))
			if !conditionMatch {
				return false
			}
		}
		return true
	})
	//isNotJustWhitespace validator to check that field not contain only white space characters only
	validate.AddValidator("isNotJustWhitespace", func(val any) bool {
		return len(strings.TrimSpace(val.(string))) > 0
	})
	validate.AddGlobalMessages(map[string]string{
		"passwordRegex":       "{field} atleast contain one small letter, one capital letter, one digit and one special character.",
		"isNotJustWhitespace": "{field} has only whitspace characters.",
	})
}

// func (f UserForm) Messages() map[string]string {
// 	return validate.MS{
// 		"required": "oh! the {field} is required",
// 		"email": "email is invalid",
// 		"Name.required": "message for special field",
// 		"Age.int": "age must int",
// 		"Age.min": "age min value is 1",
// 	}
// }

type CustomError struct {
	StatusCode   int           `json:"statusCode" validate:"required" `
	ErrorMessage string        `json:"errorMessage" validate:"required" `
	InvalidData  []InvalidData `json:"invalidData" validate:"omitempty" `
}

type InvalidData struct {
	Field string            `json:"field" `
	Error map[string]string `json:"error" `
}

func (c CustomError) Error() string {
	return c.ErrorMessage
}

func CreateCustomError(errorMessage string, statusCode int, invalidData ...InvalidData) *gqlerror.Error {
	return &gqlerror.Error{
		Message: errorMessage,
		Extensions: map[string]interface{}{
			"StatusCode":  statusCode,
			"InvalidData": invalidData,
		},
	}
}

var (
	MarshalError            = CreateCustomError("Error while marshling data.", http.StatusInternalServerError)
	UnmarshalError          = CreateCustomError("Error while unmarshling data.", http.StatusBadRequest)
	InternalServerError     = CreateCustomError("Internal Server Error.", http.StatusInternalServerError)
	BcryptError             = CreateCustomError("Error at bcypting.", http.StatusInternalServerError)
	UserAlreadyExist        = CreateCustomError("User already exist.", http.StatusConflict)
	UserDoesNotExist        = CreateCustomError("User does not exist.", http.StatusNotFound)
	MovieDoesNotExist       = CreateCustomError("Movie does not exist.", http.StatusNotFound)
	MovieReviewDoesNotExist = CreateCustomError("Movie review does not exist.", http.StatusNotFound)
	MovieReviewAlreadyExist = CreateCustomError("Movie review already exist.", http.StatusConflict)
	MovieTitleAlreadyExist  = CreateCustomError("Movie title already exist.", http.StatusConflict)
	HeaderDataMissing       = CreateCustomError("Required header not found.", http.StatusBadRequest)
	InvalidDetails          = CreateCustomError("Invalid details provided.", http.StatusBadRequest)
	JWTErrSignatureInvalid  = CreateCustomError("Invalid signature on jwt token.", http.StatusUnauthorized)
	JWTTokenInvalid         = CreateCustomError("Invalid jwt token.", http.StatusBadRequest)
	JWTTokenInvalidDetails  = CreateCustomError("Invalid jwt token details.", http.StatusBadRequest)
	JWTTokenGenerateError   = CreateCustomError("Error at generating jwt token.", http.StatusInternalServerError)
	AdminAccessRights       = CreateCustomError("Only admin have permission.", http.StatusForbidden)
	ReviewLimitExceeded     = CreateCustomError("You have exceeded the maximum number of requests allowed. Please try again later.", http.StatusTooManyRequests)

	NotNullConstraintError    = CreateCustomError("Required field cannot be empty or null. Please provide a value for the field.", http.StatusBadRequest)
	ForeignKeyConstraintError = CreateCustomError("Data doesn't exist.", http.StatusConflict)
	UniqueKeyConstraintError  = CreateCustomError("Data already exists.", http.StatusConflict)
	CheckConstraintError      = CreateCustomError("Data doesn't meet the required criteria.", http.StatusBadRequest)
	NoRowsError               = CreateCustomError("Data doesn't exist.", http.StatusNotFound)
	NoRowsAffectedError       = CreateCustomError("No data change.", http.StatusNotFound)
)

func DatabaseErrorHandling(err error) error {
	if dbErr, ok := err.(*pq.Error); ok {
		errCode := dbErr.Code
		switch errCode {
		case "23502":
			// not-null constraint violation
			return NotNullConstraintError

		case "23503":
			// foreign key violation
			return ForeignKeyConstraintError

		case "23505":
			// unique constraint violation
			return UniqueKeyConstraintError

		case "23514":
			// check constraint violation
			return CheckConstraintError

		}

	}
	return InternalServerError
}

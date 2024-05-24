package utils

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	error_handling "movie-review/error"

	"github.com/gookit/validate"
)

func ValidateStruct(data interface{}, addValidationRules map[string]string) error {
	var errorMessage []string
	validator := validate.Struct(data)
	validate.AddValidator("passwordRegex", func(val any) bool {
		password, _ := validator.Get("Password")
		conditions := []string{".{7,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
		for _, condition := range conditions {
			conditionMatch, _ := regexp.MatchString(condition, password.(string))
			if !conditionMatch {
				return false
			}
		}
		return true
	})
	validator.StringRules(addValidationRules)
	if !(validator.Validate()) {
		var invalidDataArray []error_handling.InvalidData
		errors := validator.Errors.All()
		fmt.Println(errors) // all error messages
		for key, value := range errors {
			invalidData := error_handling.InvalidData{
				Field: key,
				Error: value,
			}
			invalidDataArray = append(invalidDataArray, invalidData)
			errorMessage = append(errorMessage, key)
		}
		return error_handling.CreateCustomError("Error in field:"+strings.Join(errorMessage, ","), http.StatusBadRequest, invalidDataArray...)
	}
	return nil
}

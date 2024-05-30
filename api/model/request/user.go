package request

type UserLogin struct {
	Email    string `json:"email" validate:"required|email|max_len:320|isNotJustWhitespace"`
	Password string `json:"password" validate:"required|min_len:8|passwordRegex|isNotJustWhitespace"`
}

type UserSignup struct {
	Email     string `json:"email" validate:"required|email|max_len:320|isNotJustWhitespace"`
	Password  string `json:"password" validate:"required|min_len:8|passwordRegex|isNotJustWhitespace"`
	FirstName string `json:"firstName" validate:"required|min_len:2|max_len:255|isNotJustWhitespace"`
	LastName  string `json:"lastName" validate:"required|min_len:2|max_len:255|isNotJustWhitespace"`
}

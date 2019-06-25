package models

//easyjson:json
type UserProfile struct {
	UserID   uint   `json:"id" db:"user_id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty" valid:"stringlength(4|32)~Password must be at least 4 characters and no more than 32 characters"`
	Email    string `json:"email,omitempty" valid:"required~Email can not be empty,email~Invalid email"`
}

//easyjson:json
type ProfileError struct {
	Field string `json:"field" example:"nickname"`
	Text  string `json:"text" example:"This nickname is already taken."`
}

//easyjson:json
type ProfileErrorList struct {
	Errors []ProfileError `json:"error"`
}

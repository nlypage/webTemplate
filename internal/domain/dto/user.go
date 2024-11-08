package dto

// UserRegister @Description User registration dto
type UserRegister struct {
	Email    string `json:"email" validate:"required,email" example:"example@gmail.com"`  // Required, email must be valid
	Password string `json:"password" validate:"required,password" example:"Password1234"` // Required, password must meet certain requirements: must has upper case letters, lower case letters and digits
	Username string `json:"username" validate:"required,username" example:"linuxflight"`  // Required, user's username
}

type UserReturn struct {
	ID            string `json:"id" example:"123"`                  // User ID
	Email         string `json:"email" example:"example@gmail.com"` // User's email
	VerifiedEmail bool   `json:"verified_email" example:"true"`     // Boll variable showing, whether user's email is verified or not
	Username      string `json:"username" example:"linuxflight"`    // User's username
	Role          string `json:"role" example:"manager"`            // User's role (e.g. "Client", "Manager" etc)
}

type UserRegisterResponse struct {
	User   UserReturn `json:"user"`   // User object
	Tokens AuthTokens `json:"tokens"` // Two JWT tokens: Access token and Refresh token
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email" example:"example@gmail.com"`  // User's email, must be valid email address
	Password string `json:"password" validate:"required,password" example:"Password1234"` // User's password
}

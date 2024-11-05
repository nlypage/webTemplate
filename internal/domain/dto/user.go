package dto

type UserRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Username string `json:"username" validate:"required,username"`
}

type UserReturn struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Username      string `json:"username"`
	Role          string `json:"role"`
}

type UserRegisterResponse struct {
	User   UserReturn `json:"user"`
	Tokens AuthTokens `json:"tokens"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

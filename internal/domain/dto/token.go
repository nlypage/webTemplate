package dto

import "time"

type Token struct {
	Token   string    `json:"token" example:"somelong.token.string"`            // Token string itself
	Expires time.Time `json:"expires" example:"2024-12-08T10:00:12.961568771Z"` // Token expiration time in ISO 8601 format
}

type AuthTokens struct {
	Access  Token `json:"access"`  // Access token
	Refresh Token `json:"refresh"` // Refresh token
}

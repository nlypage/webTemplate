package dto

import "time"

type Token struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type AuthTokens struct {
	Access  Token `json:"access"`
	Refresh Token `json:"refresh"`
}

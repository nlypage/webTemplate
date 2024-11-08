package dto

type HTTPError struct {
	Code    int    `json:"code" example:"400"`               // HTTP error code
	Message string `json:"message" example:"you are retard"` // Error message
}

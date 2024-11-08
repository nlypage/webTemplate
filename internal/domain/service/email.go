package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"webTemplate/internal/adapters/config"
)

type EmailApi interface {
	Send(ctx context.Context, email string, text string, subject string) error
	Check(ctx context.Context, email string) (bool, error)
}

type emailService struct {
	service EmailApi
	config  config.MailerooConfig
}

type sendResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type checkResponse struct {
	Success   bool   `json:"success"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	Data      struct {
		Email            string `json:"email"`
		FormatValid      bool   `json:"format_valid"`
		MxFound          bool   `json:"mx_found"`
		Disposable       bool   `json:"disposable"`
		Role             bool   `json:"role"`
		Free             bool   `json:"free"`
		DomainSuggestion string `json:"domain_suggestion"`
	} `json:"data"`
}

func NewEmailService(config config.MailerooConfig) *emailService {
	return &emailService{
		config: config,
	}
}

// Send is a method to send email using https://maileroo.com API
func (s *emailService) Send(ctx context.Context, email string, text string, subject string) error {
	// request payload
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("from", s.config.FromEmail)
	_ = writer.WriteField("to", email)
	_ = writer.WriteField("subject", subject)
	_ = writer.WriteField("html", text)
	_ = writer.Close()

	// send http request
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://smtp.maileroo.com/send", payload)
	req.Header.Set("X-API-Key", s.config.SendingApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, respErr := client.Do(req)
	if respErr != nil {
		return respErr
	}

	// read body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, ioErr := io.ReadAll(res.Body)
	if ioErr != nil {
		return ioErr
	}
	var result sendResponse
	if jsonErr := json.Unmarshal(body, &result); jsonErr != nil {
		return jsonErr
	}

	// check if response is successful
	if !result.Success {
		return errors.New(result.Message)
	}

	return nil
}

// Check is a method to check via https://maileroo.com API the email address of a user
func (s *emailService) Check(ctx context.Context, email string) (bool, error) {
	// send http api request
	requestData := map[string]string{
		"email_address": email,
	}
	jsonValue, _ := json.Marshal(requestData)
	req, _ := http.NewRequest("POST", "https://verify.maileroo.net/check", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.config.VerificationApiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}

	// read body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	body, ioErr := io.ReadAll(response.Body)
	if ioErr != nil {
		return false, ioErr
	}
	var result checkResponse
	if jsonErr := json.Unmarshal(body, &result); jsonErr != nil {
		return false, jsonErr
	}

	// check if response is successful
	if !result.Success {
		return false, fmt.Errorf("%s - %s", result.ErrorCode, result.Message)
	}

	// if email is temp, or invalid format, or not found via MX record
	if !result.Data.FormatValid || !result.Data.MxFound || result.Data.Disposable {
		return false, nil
	}

	return true, nil
}

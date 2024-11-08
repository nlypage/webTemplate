package v1

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"time"
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/controller/api/validator"
	"webTemplate/internal/adapters/database/postgres"
	"webTemplate/internal/adapters/logger"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/entity"
	"webTemplate/internal/domain/service"
	"webTemplate/internal/domain/utils/auth"
)

type UserService interface {
	Create(ctx context.Context, registerReq dto.UserRegister, code string) (*entity.User, error)
	GetByID(ctx context.Context, uuid string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type TokenService interface {
	GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error)
	GenerateToken(ctx context.Context, userID string, expires time.Time, tokenType string) (*entity.Token, error)
}

type EmailService interface {
	Send(ctx context.Context, email string, text string, subject string) error
	Check(ctx context.Context, email string) (bool, error)
}

type UserHandler struct {
	userService  UserService
	tokenService TokenService
	emailService EmailService
	validator    *validator.Validator
}

func NewUserHandler(app *app.App) *UserHandler {
	userStorage := postgres.NewUserStorage(app.DB)
	tokenStorage := postgres.NewTokenStorage(app.DB)

	return &UserHandler{
		userService:  service.NewUserService(userStorage),
		tokenService: service.NewTokenService(tokenStorage),
		emailService: service.NewEmailService(app.Maileroo),
		validator:    app.Validator,
	}
}

// register godoc
// @Summary      Register a new user
// @Description  Register a new user using his email, username and password. Returns his ID, email, username, verifiedEmail boolean variable and role
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body body  dto.UserRegister true  "User registration body object"
// @Success      201  {object}  dto.UserRegisterResponse
// @Failure      400  {object}  dto.HTTPError
// @Failure      500  {object}  dto.HTTPError
// @Router       /user/register [post]
func (h UserHandler) register(c *fiber.Ctx) error {
	var userDTO dto.UserRegister

	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errValidate := h.validator.ValidateData(userDTO); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: errValidate.Error(),
		})
	}

	mailValid, mvErr := h.emailService.Check(c.Context(), userDTO.Email)
	if mvErr != nil || !mailValid {
		logger.Log.Errorf("invalid email: %s", userDTO.Email)
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "invalid email",
		})
	}

	code := auth.GenerateCode()
	msErr := h.emailService.Send(c.Context(), userDTO.Email, fmt.Sprintf("Your code is: <b>%s</b>", code), "Verification Code")
	if msErr != nil {
		logger.Log.Errorf("email sending error: %s", msErr.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: msErr.Error(),
		})
	}

	user, errCreate := h.userService.Create(c.Context(), userDTO, code)
	if errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: errCreate.Error(),
		})
	}

	tokens, tokensErr := h.tokenService.GenerateAuthTokens(c.Context(), user.ID)
	if tokensErr != nil || tokens == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to generate auth tokens",
		})
	}

	response := dto.UserRegisterResponse{
		User: dto.UserReturn{
			ID:            user.ID,
			Email:         user.Email,
			VerifiedEmail: user.VerifiedEmail,
			Username:      user.Username,
			Role:          user.Role,
		},
		Tokens: *tokens,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// login godoc
// @Summary      Login to existing user account.
// @Description  Login to existing user account using his email, username and password. Returns his ID, email, username, verifiedEmail boolean variable and role
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body body  dto.UserLogin true  "User login body object"
// @Success      200  {object}  dto.UserRegisterResponse
// @Failure      400  {object}  dto.HTTPError
// @Failure      403  {object}  dto.HTTPError
// @Failure      404  {object}  dto.HTTPError
// @Failure      500  {object}  dto.HTTPError
// @Router       /user/login [post]
func (h UserHandler) login(c *fiber.Ctx) error {
	var userDTO dto.UserLogin

	if err := c.BodyParser(&userDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errValidate := h.validator.ValidateData(userDTO); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: errValidate.Error(),
		})
	}

	user, errFetch := h.userService.GetByEmail(c.Context(), userDTO.Email)
	if errFetch != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.HTTPError{
			Code:    fiber.StatusNotFound,
			Message: "not found",
		})
	}

	passErr := user.ComparePassword(userDTO.Password)
	if passErr != nil {
		return c.Status(fiber.StatusForbidden).JSON(dto.HTTPError{
			Code:    fiber.StatusForbidden,
			Message: "invalid password",
		})
	}

	tokens, tokensErr := h.tokenService.GenerateAuthTokens(c.Context(), user.ID)
	if tokensErr != nil || tokens == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: "failed to generate auth tokens",
		})
	}

	response := dto.UserRegisterResponse{
		User: dto.UserReturn{
			ID:            user.ID,
			Email:         user.Email,
			VerifiedEmail: user.VerifiedEmail,
			Username:      user.Username,
			Role:          user.Role,
		},
		Tokens: *tokens,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// refreshToken godoc
// @Summary      Refresh the access token
// @Description  Get a new access token using a valid refresh token
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body body  dto.Token true  "Access token object"
// @Success      200  {object}  dto.Token
// @Failure      400  {object}  dto.HTTPError
// @Failure      403  {object}  dto.HTTPError
// @Failure      500  {object}  dto.HTTPError
// @Router       /user/refresh [post]
func (h UserHandler) refreshToken(c *fiber.Ctx) error {
	var accessTokenDTO dto.Token

	if err := c.BodyParser(&accessTokenDTO); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errValidate := h.validator.ValidateData(accessTokenDTO); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: errValidate.Error(),
		})
	}

	userID, errToken := auth.VerifyToken(accessTokenDTO.Token, viper.GetString("service.backend.jwt.secret"), auth.TokenTypeAccess)

	if errToken != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Code:    fiber.StatusUnauthorized,
			Message: errToken.Error(),
		})
	}

	expTime := time.Now().UTC().Add(time.Minute * time.Duration(viper.GetInt("service.backend.jwt.access-token-expiration")))

	newAccess, errNewAccess := h.tokenService.GenerateToken(c.Context(),
		userID,
		expTime,
		auth.TokenTypeAccess)

	if errNewAccess != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: errNewAccess.Error(),
		})
	}

	response := dto.Token{
		Token:   newAccess.Token,
		Expires: expTime,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// verify godoc
// @Summary      Verify user account
// @Description  Verify a user account with a code, sent to user's email
// @Tags         user
// @Accept       json
// @Produce      json
// @Security Bearer
// @Param        body body  dto.UserCode true  "User's email code"
// @Success      200  {object}  dto.HTTPStatus
// @Failure      400  {object}  dto.HTTPError
// @Failure      401  {object}  dto.HTTPError
// @Failure      403  {object}  dto.HTTPError
// @Failure      500  {object}  dto.HTTPError
// @Router       /user/verify [post]
func (h UserHandler) verify(c *fiber.Ctx) error {
	var userCode dto.UserCode

	if err := c.BodyParser(&userCode); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errValidate := h.validator.ValidateData(userCode); errValidate != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: errValidate.Error(),
		})
	}

	user, authErr := auth.GetUserFromJWT(c.Get("Authorization"), auth.TokenTypeAccess, c.Context(), h.userService.GetByID)

	if authErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.HTTPError{
			Code:    fiber.StatusUnauthorized,
			Message: authErr.Error(),
		})
	}

	if user.VerificationCode == "NULL" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.HTTPError{
			Code:    fiber.StatusBadRequest,
			Message: "already verified",
		})
	}

	if user.VerificationCode != userCode.Code {
		return c.Status(fiber.StatusForbidden).JSON(dto.HTTPError{
			Code:    fiber.StatusForbidden,
			Message: "invalid code",
		})
	}

	user.VerificationCode = "NULL"
	_, updateErr := h.userService.Update(c.Context(), user)
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.HTTPError{
			Code:    fiber.StatusInternalServerError,
			Message: updateErr.Error(),
		})
	}

	response := dto.HTTPStatus{
		Code:    200,
		Message: "email verified",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h UserHandler) Setup(router fiber.Router, middleware fiber.Handler) {
	userGroup := router.Group("/user")
	userGroup.Post("/register", h.register)
	userGroup.Post("/login", h.login)
	userGroup.Post("/refresh", h.refreshToken)
	userGroup.Post("/verify", h.verify, middleware)
}

package v1

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/controller/api/validator"
	"webTemplate/internal/adapters/database/postgres"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/entity"
	"webTemplate/internal/domain/service"
)

type UserService interface {
	Create(ctx context.Context, registerReq dto.UserRegister) (*entity.User, error)
}

type TokenService interface {
	GenerateAuthTokens(c context.Context, userID string) (*dto.AuthTokens, error)
}

type UserHandler struct {
	userService  UserService
	tokenService TokenService
	validator    *validator.Validator
}

func NewUserHandler(app *app.App) *UserHandler {
	userStorage := postgres.NewUserStorage(app.DB)
	tokenStorage := postgres.NewTokenStorage(app.DB)

	return &UserHandler{
		userService:  service.NewUserService(userStorage),
		tokenService: service.NewTokenService(tokenStorage),
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

	user, errCreate := h.userService.Create(c.Context(), userDTO)
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
			Message: tokensErr.Error(),
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

func (h UserHandler) Setup(router fiber.Router) {
	userGroup := router.Group("/user")
	userGroup.Post("/register", h.register)
}

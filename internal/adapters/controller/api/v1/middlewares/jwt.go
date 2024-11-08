package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/database/postgres"
	"webTemplate/internal/domain/common/errorz"
	"webTemplate/internal/domain/entity"
	"webTemplate/internal/domain/service"
	"webTemplate/internal/domain/utils/auth"
)

type UserService interface {
	GetByID(ctx context.Context, uuid string) (*entity.User, error)
}

type MiddlewareHandler struct {
	userService UserService
}

// NewMiddlewareHandler is a function that returns a new instance of MiddlewareHandler.
func NewMiddlewareHandler(app *app.App) *MiddlewareHandler {
	userStorage := postgres.NewUserStorage(app.DB)
	userService := service.NewUserService(userStorage)

	return &MiddlewareHandler{
		userService: userService,
	}
}

// IsAuthenticated is a function that checks whether the user has sufficient rights to access the endpoint
/*
 * c *fiber.Ctx - request context
 * role string - the user role required for access
 * tokenType string - the type of token that is required to access the endpoint
 */
func (h MiddlewareHandler) IsAuthenticated(c *fiber.Ctx, role string, tokenType string) error {
	if len(c.GetReqHeaders()["Authorization"]) == 0 {
		return errorz.AuthHeaderIsEmpty
	}

	authHeader := c.GetReqHeaders()["Authorization"][0]
	if authHeader == "" {
		return errorz.AuthHeaderIsEmpty
	}

	id, errVerify := auth.VerifyToken(authHeader, viper.GetString("service.backend.jwt.secret"), tokenType)
	if errVerify != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errVerify.Error(),
		})
	}

	user, errGetUser := h.userService.GetByID(c.Context(), id)
	if errGetUser != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": errGetUser.Error(),
		})
	}

	if user.Role != role {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": errorz.Forbidden.Error(),
		})
	}
	return c.Next()

}

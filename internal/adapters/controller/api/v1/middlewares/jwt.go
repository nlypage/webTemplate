package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/config"
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
 * tokenType string - the type of token that is required to access the endpoint
 * requiredRights ...string - the rights that the user must have
 */
func (h MiddlewareHandler) IsAuthenticated(tokenType string, requiredRights ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

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

		if !config.RoleHasRights(user.Role, requiredRights) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": errorz.Forbidden,
			})
		}
		return c.Next()
	}
}

package middlewares

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strings"
	"webTemplate/cmd/app"
	"webTemplate/internal/adapters/database/postgres"
	"webTemplate/internal/domain/entity"
	"webTemplate/internal/domain/service"
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

func (h MiddlewareHandler) IsAuthenticated(c *fiber.Ctx, role string, tokenType string) error {
	if len(c.GetReqHeaders()["Authorization"]) == 0 {
		return
	}

	authHeader := c.GetReqHeaders()["Authorization"][0]
	if authHeader == "" {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status":  false,
			"message": "auth header is empty",
		})
	}

	uuid, password, errParse := utils.ParseJwt(parts[1])
	if errParse != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errParse.Error(),
		})
	}

	user, errGetUser := h.userService.GetByUUID(c.Context(), uuid)
	if errGetUser != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errGetUser.Error(),
		})
	}

	if string(user.Password) != password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"body":   errroz.TokenExpired.Error(),
		})
	}
	return c.Next()

}
}

package v1

import (
	"github.com/gofiber/fiber/v2"
	"webTemplate/internal/domain/dto"
	"webTemplate/internal/domain/service"
)

type UserHandler struct {
	userService service.UserService
}

func (h UserHandler) register(c *fiber.Ctx) error {
	var userDTO dto.UserRegister

	if err := c.BodyParser(&userDTO); err != nil {
		return err
	}

	return nil
}

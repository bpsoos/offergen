package endpoint

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (e *Handler) Profile(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
}

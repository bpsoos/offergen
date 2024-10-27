package preview

import (
	"github.com/gofiber/fiber/v2"
)

func (i *Handler) Delete(ctx *fiber.Ctx) error {
	return ctx.SendStatus(200)
}

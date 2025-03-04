package inventory

import (
	"offergen/endpoint/handlermachinery"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) BatchGetCategory(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	categories, err := h.inventoryManager.BatchGetCountedCategory(
		userID,
	)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.inventoryTemplater.Categories(c, w,
		userID,
		categories,
	)
}

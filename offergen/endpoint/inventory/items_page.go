package inventory

import "github.com/gofiber/fiber/v2"

func (h *Handler) ItemsPage(ctx *fiber.Ctx) error {
	return h.renderer.Render(
		ctx,
		h.inventoryTemplater.ItemsPage(h.tokenVerifier.GetUserID(ctx)),
	)
}

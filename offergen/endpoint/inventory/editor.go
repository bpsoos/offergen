package inventory

import "github.com/gofiber/fiber/v2"

func (h *Handler) Editor(ctx *fiber.Ctx) error {
	return h.renderer.Render(
		ctx,
		h.pageTemplater.Inventory(h.tokenVerifier.GetUserID(ctx)),
	)
}

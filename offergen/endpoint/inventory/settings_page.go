package inventory

import "github.com/gofiber/fiber/v2"

func (h *Handler) SettingsPage(ctx *fiber.Ctx) error {
	inv, err := h.inventoryManager.GetInventory(h.tokenVerifier.GetUserID(ctx))

	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return h.renderer.Render(ctx, h.inventoryTemplater.SettingsPage(inv))
}

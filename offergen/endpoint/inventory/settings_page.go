package inventory

import "github.com/gofiber/fiber/v2"

func (h *Handler) SettingsPage(ctx *fiber.Ctx) error {
	return h.renderer.Render(ctx, h.inventoryTemplater.SettingsPage())
}

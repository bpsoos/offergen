package inventory

import "github.com/gofiber/fiber/v2"

func (h *Handler) Editor(ctx *fiber.Ctx) error {
	items, err := h.inventoryManager.BatchGetItem(0, 10, h.tokenVerifier.GetUserID(ctx))
	if err != nil {
		logger.Error("error fetching items", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return h.renderer.Render(ctx, h.pageTemplater.Inventory(items))
}

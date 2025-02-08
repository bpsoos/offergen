package offerings

import "github.com/gofiber/fiber/v2"

func (h *Handler) GetOffering(ctx *fiber.Ctx) error {
	ownerID := ctx.Params("owner_id", "")
	if ownerID == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	items, err := h.inventoryManager.BatchGetItem(0, 1000, ownerID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return h.renderer.Render(ctx, h.offeringTemplater.Menu(items))
}

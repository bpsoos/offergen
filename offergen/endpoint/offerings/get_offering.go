package offerings

import (
	"offergen/endpoint/handlermachinery"
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetOffering(ctx *fiber.Ctx) error {
	ownerID := ctx.Params("owner_id", "")
	if ownerID == "" {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	inv, err := h.inventoryManager.GetInventory(ownerID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	if !inv.IsPublished {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	items, err := h.inventoryManager.BatchGetItem(
		ownerID,
		&models.GetItemsInput{From: 0, Amount: 100},
	)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.offeringTemplater.Offering(c, w, inv.Title, items)
}

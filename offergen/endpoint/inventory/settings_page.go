package inventory

import (
	"offergen/endpoint/handlermachinery"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SettingsPage(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	logger.Info("getting inventory", "user", userID)
	inv, err := h.inventoryManager.GetInventory(userID)

	if err != nil {
		logger.Error("error getting inventory", "err", err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.inventoryTemplater.SettingsPage(c, w,
		userID,
		inv,
	)
}

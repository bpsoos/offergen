package inventory

import (
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	input := models.DeleteItemInput{ItemID: ctx.Params("id")}
	if err := h.structValidator.Validate(input); err != nil {
		logger.Error("delete item missing item id", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	logger.Info("deleting item", "userID", userID, "itemID", input.ItemID)
	err := h.inventoryManager.DeleteItem(input.ItemID, userID)
	if err != nil {
		logger.Error("error deleting item", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	ctx.Response().Header.Add("HX-Trigger", "item-deleted")

	return ctx.SendStatus(fiber.StatusOK)
}

package inventory

import (
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Items(ctx *fiber.Ctx) error {
	input := new(models.GetItemsInput)
	if err := ctx.QueryParser(input); err != nil {
		logger.Error("could not parse get items query", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := h.structValidator.Validate(input); err != nil {
		logger.Error("invalid get items query", "errMsg", err.Error(), "amount", input.Amount, "from", input.From)
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	items, err := h.inventoryManager.BatchGetItem(
		input.From,
		input.Amount,
		h.tokenVerifier.GetUserID(ctx),
	)
	if err != nil {
		logger.Error("error fetching items", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	logger.Info("fetched items", "itemCount", len(items))
	return h.renderer.Render(ctx, h.inventoryTemplater.Items(items))
}

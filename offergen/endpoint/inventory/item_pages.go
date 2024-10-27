package inventory

import (
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ItemPages(ctx *fiber.Ctx) error {
	input := new(models.GetItemPagesInput)
	if err := ctx.QueryParser(input); err != nil {
		logger.Error("could not parse get item pages query", "errMsg", err.Error())
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	if err := h.structValidator.Validate(input); err != nil {
		logger.Error("invalid get item pages query", "errMsg", err.Error(), "current", input.Current)
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	userID := h.tokenVerifier.GetUserID(ctx)
	count, err := h.inventoryManager.ItemCount(userID)
	if err != nil {
		logger.Error("error getting item count", "userID", userID, "errMsg", err)
	}

	var lastPage int
	if count == 0 {
		lastPage = 0
	} else {
		lastPage = (count-1)/10 + 1
	}
	logger.Info("last page", "lastPage", lastPage, "count", count)
	return h.renderer.Render(ctx, h.inventoryTemplater.Paginator(int(input.Current), lastPage))
}

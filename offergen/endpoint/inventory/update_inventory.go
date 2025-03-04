package inventory

import (
	"offergen/endpoint"
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdateInventory(ctx *fiber.Ctx) error {
	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if values == nil {
		logger.Info("no args in request")
		return h.renderUpdateInventoryError(ctx, "missing fields")
	}

	input := new(models.UpdateInventoryInput)
	err := h.formDecoder.Decode(input, values)
	if err != nil {
		errDecode := h.formDecoder.MustParseDecodeErrors(err)
		logger.Info("decode error", "errMsg", errDecode[0].Error())

		return h.renderUpdateInventoryError(ctx, errDecode[0].Field()+": validation error")
	}
	logger.Info(
		"decoded",
		"title", input.Title,
		"isPublished", input.IsPublished,
	)

	err = h.structValidator.Validate(input)
	if err != nil {
		validationErrors := h.structValidator.MustParseValidationErrors(err)
		logger.Info("valdation error", "errMsg", validationErrors[0].Error())

		return h.renderUpdateInventoryError(ctx, validationErrors[0].Field()+": validation error")
	}
	userID := h.tokenVerifier.GetUserID(ctx)

	logger.Info(
		"decoded and parsed update inventory input",
		"title", input.Title,
		"isPublished", input.IsPublished,
		"userID", userID,
	)

	inv, err := h.inventoryManager.UpdateInventory(
		userID,
		&models.UpdateInventoryInput{
			Title:       input.Title,
			IsPublished: input.IsPublished,
		},
	)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	logger.Info(
		"updated",
		"title", inv.Title,
		"isPublished", inv.IsPublished,
	)

	return h.inventoryTemplater.InventoryDetails(
		ctx.Context(),
		ctx.Response().BodyWriter(),
		inv,
	)
}

func (h *Handler) renderUpdateInventoryError(ctx *fiber.Ctx, errorMsg string) error {
	return ctx.SendStatus(fiber.StatusUnprocessableEntity)
}

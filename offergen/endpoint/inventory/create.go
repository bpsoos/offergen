package inventory

import (
	"offergen/endpoint"
	"offergen/endpoint/models"
	"regexp"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Create(ctx *fiber.Ctx) error {
	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if values == nil {
		logger.Info("no args in request")
		return h.renderItemCreateError(ctx, "missing fields")
	}

	input := new(models.AddItemInput)
	err := h.formDecoder.Decode(input, values)
	if err != nil {
		errDecode := h.formDecoder.MustParseDecodeErrors(err)
		logger.Info("decode error", "errMsg", errDecode[0].Error())

		return h.renderItemCreateError(ctx, errDecode[0].Field()+": validation error")
	}

	err = h.structValidator.Validate(input)
	if err != nil {
		validationErrors := h.structValidator.MustParseValidationErrors(err)
		logger.Info("valdation error", "errMsg", validationErrors[0].Error())

		return h.renderItemCreateError(ctx, validationErrors[0].Field()+": validation error")
	}

	pattern, err := regexp.Compile(models.AllowedNamePattern)
	if err != nil {
		logger.Error(
			"allowed name pattern is an invalid regex",
			"allowedNamePattern", models.AllowedNamePattern,
		)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	if !pattern.MatchString(input.Name) {
		logger.Info("allow pattern validation error", "inputName", input.Name)
		return h.renderItemCreateError(ctx, "Name: validation error")
	}

	logger.Info(
		"parsed and validated item",
		"itemName", input.Name,
		"itemPrice", input.Price,
	)
	_, err = h.inventoryManager.CreateItem(input, h.tokenVerifier.GetUserID(ctx))
	if err != nil {
		logger.Info("could not create item", "inputName", input.Name)
		return h.renderItemCreateError(ctx, "something went wrong")
	}

	return h.renderer.Render(ctx, h.pageTemplater.Inventory(nil))
}

func (h *Handler) CreatePage(ctx *fiber.Ctx) error {
	return h.renderer.Render(ctx, h.inventoryTemplater.ItemCreator())
}

func (h *Handler) renderItemCreateError(ctx *fiber.Ctx, message string) error {
	ctx.Response().Header.Add("HX-Retarget", ".error-msg")

	return h.renderer.Render(ctx, h.errorTemplater.ErrorMessage(message))
}

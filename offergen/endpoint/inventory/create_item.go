package inventory

import (
	"encoding/json"
	"offergen/endpoint"
	"offergen/endpoint/handlermachinery"
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateItem(ctx *fiber.Ctx) error {
	createItemInput, err := h.parseItem(ctx)
	if err != nil {
		return err
	}

	err = h.structValidator.Validate(createItemInput)
	if err != nil {
		validationErrors := h.structValidator.MustParseValidationErrors(err)
		logger.Info("valdation error", "errMsg", validationErrors[0].Error())

		return h.renderItemCreateError(ctx, validationErrors[0].Field()+": validation error")
	}

	logger.Info(
		"parsed and validated item",
		"itemName", createItemInput.Name,
		"itemPrice", createItemInput.Price,
		"desc", createItemInput.Desc,
	)

	item, err := h.inventoryManager.CreateItem(createItemInput, h.tokenVerifier.GetUserID(ctx))
	if err != nil {
		logger.Info("could not create item", "inputName", createItemInput.Name)
		return h.renderItemCreateError(ctx, "something went wrong")
	}

	if acceptsJSON(ctx) {
		return ctx.JSON(&item)
	}

	return h.inventoryTemplater.Inventory(
		ctx.Context(),
		ctx.Response().BodyWriter(),
		h.tokenVerifier.GetUserID(ctx),
	)
}

func acceptsJSON(ctx *fiber.Ctx) bool {
	return string(ctx.Request().Header.Peek(fiber.HeaderAccept)) == fiber.MIMEApplicationJSON
}

func isJSONContentType(ctx *fiber.Ctx) bool {
	return string(ctx.Request().Header.ContentType()) == fiber.MIMEApplicationJSON
}

func (h *Handler) CreatePage(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.inventoryTemplater.ItemCreator(c, w, userID)
}

func (h *Handler) parseItem(ctx *fiber.Ctx) (*models.AddItemInput, error) {
	input := new(models.AddItemInput)
	if isJSONContentType(ctx) {
		body := ctx.Body()
		err := json.Unmarshal(body, &input)
		if err != nil {

			logger.Info("json parse error", "err", err)

			return nil, err
		}
		return input, nil
	}

	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if values == nil {
		logger.Info("no args in request")
		return nil, h.renderItemCreateError(ctx, "missing fields")
	}

	err := h.formDecoder.Decode(input, values)
	if err != nil {
		errDecode := h.formDecoder.MustParseDecodeErrors(err)
		logger.Info("decode error", "errMsg", errDecode[0].Error())

		return nil, h.renderItemCreateError(ctx, errDecode[0].Field()+": validation error")
	}
	return input, err
}

func (h *Handler) renderItemCreateError(ctx *fiber.Ctx, message string) error {
	ctx.Response().Header.Add("HX-Retarget", ".error-msg")

	return h.renderer.Render(ctx, h.errorTemplater.ErrorMessage(message))
}

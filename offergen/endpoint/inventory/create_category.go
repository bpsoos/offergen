package inventory

import (
	"offergen/endpoint/handlermachinery"
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateCategoryInit(ctx *fiber.Ctx) error {
	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.inventoryTemplater.CreateCategoryForm(c, w)
}

func (h *Handler) CreateCategory(ctx *fiber.Ctx) error {
	values := handlermachinery.ParseURLValues(ctx)
	if len(values) == 0 {
		logger.Info("no args in request")
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}
	input := new(models.CreateCategoryInput)
	err := h.formDecoder.Decode(input, values)
	if err != nil {
		errDecode := h.formDecoder.MustParseDecodeErrors(err)
		logger.Info("decode error", "errMsg", errDecode[0].Error())

		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	err = h.inventoryManager.CreateCategory(
		h.tokenVerifier.GetUserID(ctx),
		input.Name,
	)
	if err != nil {
		logger.Error("error creating category", "err", err)

		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Response().Header.Add("HX-Redirect", "/inventory/categories")
	return ctx.SendStatus(fiber.StatusOK)
}

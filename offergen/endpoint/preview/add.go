package preview

import (
	"offergen/endpoint"
	"offergen/endpoint/models"
	"offergen/templates"
	itemTemplates "offergen/templates/components/items"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (i *Handler) Add(ctx *fiber.Ctx) error {
	input := new(models.AddItemInput)
	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if values == nil {
		logger.Info("no args in request")
		return renderItemAddError(ctx, "missing fields")
	}

	err := i.formDecoder.Decode(input, values)
	if err != nil {
		errDecode := i.formDecoder.MustParseDecodeErrors(err)
		logger.Info("decode error", "errMsg", errDecode[0].Error())

		return renderItemAddError(ctx, errDecode[0].Field()+": validation error")
	}

	err = i.structValidator.Validate(input)
	if err != nil {
		validationErrors := i.structValidator.MustParseValidationErrors(err)
		logger.Info("valdation error", "errMsg", validationErrors[0].Error())

		return renderItemAddError(ctx, validationErrors[0].Field()+": validation error")
	}

	pattern, err := regexp.Compile(models.AllowedNamePattern)
	if err != nil {
		logger.Error(
			"allowed name pattern is an invalid regex",
			"allowedNamePattern", models.AllowedNamePattern,
		)
		panic(err)
	}

	if !pattern.MatchString(input.Name) {
		logger.Info("allow pattern validation error", "inputName", input.Name)
		return renderItemAddError(ctx, "Name: validation error")
	}

	return renderItemRowResponse(ctx, input)
}

func renderItemAddError(ctx *fiber.Ctx, errorMsg string) error {
	return templates.Render(
		ctx,
		itemTemplates.ItemAddError(errorMsg),
	)
}

func renderItemRowResponse(ctx *fiber.Ctx, input *models.AddItemInput) error {
	return templates.Render(
		ctx,
		itemTemplates.ItemRowResponse(
			models.Item{
				ID:    uuid.New(),
				Name:  input.Name,
				Price: input.Price,
			},
		),
	)
}

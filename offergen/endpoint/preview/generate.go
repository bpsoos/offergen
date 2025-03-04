package preview

import (
	"net/http"
	"net/url"
	"offergen/endpoint"
	"offergen/endpoint/models"
	"offergen/logging"

	"github.com/go-playground/form"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) Generate(ctx *fiber.Ctx) error {
	logger := logging.GetLogger()
	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if values == nil {
		logger.Info("missing request args")
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	items := h.convertItems(&values)
	for _, item := range items {
		err := h.structValidator.Validate(&item)
		if err == nil {
			continue
		}
		validationErrors := h.structValidator.MustParseValidationErrors(err)
		logger.Error(
			"invalid item in input",
			"firstErrField", validationErrors[0].Field(),
			"firstErrMsg", validationErrors[0].Error(),
		)
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	logger.Info("got raw items", "itemsLen", len(items))
	logger.Info("parsed items")

	return h.renderer.Render(
		ctx,
		h.offeringTemplater.Offering(items),
	)
}

func (i *Handler) convertItems(values *url.Values) []models.Item {
	if values == nil || len(*values) == 0 {
		return nil
	}

	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(
		func(vals []string) (interface{}, error) {
			return uuid.Parse(vals[0])
		},
		uuid.UUID{},
	)

	items := new(models.ItemsForm)
	if err := decoder.Decode(items, *values); err != nil {
		panic(err)
	}
	itemsList := make([]models.Item, len(items.Items))

	index := 0
	for k, v := range items.Items {
		itemsList[index] = models.Item{
			ID:    uuid.MustParse(k),
			Price: v.Price,
			Name:  v.Name,
		}
		index++
	}

	return itemsList
}

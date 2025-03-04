package inventory

import (
	"offergen/endpoint/handlermachinery"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Editor(ctx *fiber.Ctx) error {
	c, w := handlermachinery.ToTemplaterContext(ctx)
	return h.inventoryTemplater.Inventory(c, w,
		h.tokenVerifier.GetUserID(ctx),
	)
}

package endpoint

import (
	"offergen/templates"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (e *Handler) Index(ctx *fiber.Ctx) error {
	return templates.Render(
		ctx,
		pages.Index(
			e.verifier.IsValidToken(
				ctx.Context(),
				ctx.Request().Header.Cookie(e.authCookieName),
			),
		),
	)
}

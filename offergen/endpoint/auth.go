package endpoint

import (
	"offergen/templates"
	"offergen/templates/components"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (e *Handler) Logout(ctx *fiber.Ctx) error {
	ctx.Response().Header.DelClientCookie(e.authCookieName)

	return templates.Render(ctx, components.Sidebar(false))
}

func (e *Handler) Authenticate(ctx *fiber.Ctx) error {
	return templates.Render(ctx, pages.Authenticate(""))
}

func (e *Handler) Unauthorized(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Hx-Redirect", "/")

	return ctx.SendStatus(fiber.StatusTemporaryRedirect)
}

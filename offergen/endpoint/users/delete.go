package users

import (
	"net/http"
	"offergen/templates"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	logger.Info("parsed user id", "userID", userID)

	if err := h.userManager.Delete(
		userID,
		h.tokenVerifier.GetUserToken(ctx),
	); err != nil {
		logger.Error("error deleting user from db", "errMsg", err.Error())

		return ctx.SendStatus(http.StatusInternalServerError)
	}

	h.delClientCookie(ctx)

	return templates.Render(ctx, pages.Index(false))
}

func (h *Handler) delClientCookie(ctx *fiber.Ctx) {
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(h.authCookieName)
	cookie.SetPath("/")
	cookie.SetExpire(fasthttp.CookieExpireDelete)

	ctx.Response().Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)
}

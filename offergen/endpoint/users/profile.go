package users

import (
	"net/http"
	"offergen/templates"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Profile(ctx *fiber.Ctx) error {
	userID := h.tokenVerifier.GetUserID(ctx)
	logger.Info("parsed user id", "userID", userID)

	email, err := h.userManager.GetEmail(userID)
	if err != nil {
		logger.Error("error getting email of user", "errMsg", err.Error())
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return templates.Render(
		ctx,
		pages.Profile(email),
	)
}

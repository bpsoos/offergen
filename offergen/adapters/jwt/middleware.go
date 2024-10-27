package jwt

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const (
	userIDCtxKey = "userID"
	tokenCtxKey  = "token"
)

func (v Verifier) VerifyUser(ctx *fiber.Ctx) error {
	rawToken := ctx.Request().Header.Cookie(v.authCookieName)
	if len(rawToken) == 0 {
		logger.Error("missing auth cookie")

		return ctx.SendStatus(http.StatusUnauthorized)
	}

	token, err := v.parseToken(ctx.Context(), rawToken)
	if err != nil {
		logger.Error("user unauthenticated", "errMsg", err.Error())

		ctx.Response().Header.Set("Hx-Redirect", "/")
		return ctx.SendStatus(fiber.StatusTemporaryRedirect)
	}

	ctx.Context().SetUserValue(userIDCtxKey, token.Subject())
	ctx.Context().SetUserValue(tokenCtxKey, rawToken)

	return ctx.Next()
}

func (v Verifier) GetUserID(ctx *fiber.Ctx) string {
	userID, ok := ctx.Context().UserValue(userIDCtxKey).(string)

	if !ok {
		panic("invalid user id type")
	}

	return userID
}

func (v Verifier) GetUserToken(ctx *fiber.Ctx) []byte {
	token, ok := ctx.Context().UserValue(tokenCtxKey).([]byte)

	if !ok {
		panic("invalid token type")
	}

	return token
}

package auth

import (
	"offergen/endpoint/models"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (auth *Handler) SignUp(ctx *fiber.Ctx) error {
	email, err := auth.parseEmail(ctx)
	if err != nil {
		logger.Error("parse email error", "errMsg", err.Error())

		return renderAuthenticateInit(ctx, "something went wrong")
	}

	logger.Info(
		"successfully decoded and validated email",
		"email", email,
	)
	flowParams, err := auth.authencticator.SignUp(email)
	if err != nil {
		logger.Error("sign up error", "errMsg", err.Error())

		return renderAuthenticateInit(ctx, "something went wrong")
	}

	logger.Info(
		"successfully started sign up flow",
		"email", email,
		"flowID", flowParams.FlowID,
	)

	ctx.Response().Header.SetCookie(
		auth.createFlowCookie(flowParams),
	)

	return renderVerifyPasscode(ctx, flowParams.Email)
}

func (auth *Handler) createFlowCookie(authFlowParams *models.AuthFlowParams) *fasthttp.Cookie {
	flowCookie := fasthttp.AcquireCookie()
	flowCookie.SetKey(auth.flowCookieName)
	flowCookie.SetHTTPOnly(true)
	flowCookie.SetMaxAge(900)
	flowCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	flowCookie.SetSecure(auth.flowCookieIsSecure)
	flowCookie.SetPath(auth.flowCookiePath)
	flowCookie.SetDomain(auth.flowCookieDomain)

	flowCookie.SetValueBytes(authFlowParams.ToEncodedJson())

	return flowCookie
}

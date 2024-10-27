package auth

import (
	"errors"
	"offergen/endpoint"
	"offergen/endpoint/models"
	"offergen/templates"
	errorTemplates "offergen/templates/components/errors"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (a *Handler) Init(ctx *fiber.Ctx) error {
	logger.Info("init called", "userAgent", string(ctx.Request().Header.Peek(fiber.HeaderUserAgent)))
	email, err := a.parseEmail(ctx)
	if err != nil {
		return renderAuthenticateInitError(ctx, "email: validation error")
	}

	isUserRegistered, err := a.authencticator.IsUserRegistered(email)
	if err != nil {
		return renderAuthenticateInitError(ctx, "error communicating with authentication server")
	}

	if !isUserRegistered {
		logger.Info("no existing user found with email address", "email", email)

		return templates.Render(ctx, pages.ConfirmSignUp(email))
	}

	logger.Info("user found with email address", "email", email)
	return a.login(ctx, email)
}

func (auth *Handler) login(ctx *fiber.Ctx, email string) error {
	result, err := auth.authencticator.LoginInit(email)
	if err != nil {
		logger.Error("login by email error", "errMsg", err.Error())

		return renderAuthenticateInitError(ctx, "something went wrong")
	}

	flowCookie := auth.createFlowCookie(result)
	ctx.Response().Header.SetCookie(flowCookie)

	return renderVerifyPasscode(ctx, email)
}

func (a *Handler) parseEmail(ctx *fiber.Ctx) (string, error) {
	values := endpoint.ToURLValues(ctx.Request().PostArgs())
	if len(values) == 0 {
		return "", errors.New("can't parse email: missing request args")
	}

	emailForm := new(models.EmailForm)

	if err := a.formDecoder.Decode(emailForm, values); err != nil {
		return "", err
	}

	if err := a.structValidator.Validate(emailForm); err != nil {
		return "", err
	}

	return emailForm.Email, nil
}

func renderAuthenticateInitError(ctx *fiber.Ctx, errorMsg string) error {
	ctx.Response().Header.Set("HX-Retarget", "#authenticate-init-form .error-msg")
	return templates.Render(ctx, errorTemplates.ErrorMessage(errorMsg))
}

func renderVerifyPasscode(ctx *fiber.Ctx, email string) error {
	return templates.Render(
		ctx,
		pages.VerifyPasscode(email),
	)
}

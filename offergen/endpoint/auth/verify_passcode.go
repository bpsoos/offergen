package auth

import (
	"errors"
	"offergen/endpoint"
	"offergen/endpoint/models"
	"offergen/templates"
	errroTemplates "offergen/templates/components/errors"
	"offergen/templates/pages"

	"github.com/gofiber/fiber/v2"
)

func (a *Handler) VerifyPasscode(ctx *fiber.Ctx) error {
	logger.Info("parsing auth flow cookie")
	authFlowParams, err := a.parseAuthFlowParams(ctx)
	if err != nil {
		logger.Error("error parsing flow cookie", "errMsg", err.Error())

		return renderAuthenticateInit(ctx, "something went wrong")
	}

	logger.Info("parsing passcode")
	passcode, err := a.parsePasscode(ctx)
	if err != nil {
		logger.Error("error getting passcode", "errMsg", err.Error())
		return renderVerifyPasscodeError(ctx, "passcode: validation error")
	}

	logger.Info("verifying passcode", "email", authFlowParams.Email, "flowID", authFlowParams.FlowID)
	authCookie, err := a.authencticator.VerifyPasscode(&models.VerifyPasscodeInput{
		AuthFlowParams: authFlowParams,
		PasscodeForm:   passcode,
	})
	if err != nil {
		if errInvalid, ok := err.(models.ErrInvalidPasscode); ok {
			logger.Info("invalid passcode received", "errMsg", err.Error())
			authFlowParams.CsrfToken = errInvalid.CsrfToken
			ctx.Response().Header.SetCookie(a.createFlowCookie(authFlowParams))

			return renderVerifyPasscodeError(ctx, "invalid passcode")
		}
		logger.Error("verify passcode error", "errMsg", err.Error())

		return renderVerifyPasscodeError(ctx, "something went wrong")
	}

	logger.Info("successfully verified passcode", "email", authFlowParams.Email, "flowID", authFlowParams.FlowID)
	ctx.Response().Header.SetCookie(authCookie)

	ctx.Response().Header.Set("HX-Retarget", "body")
	return templates.Render(ctx, pages.Index(true))
}

func (a *Handler) parseAuthFlowParams(ctx *fiber.Ctx) (*models.AuthFlowParams, error) {
	flowCookieContents := ctx.Request().Header.Cookie(a.flowCookieName)
	if len(flowCookieContents) == 0 {
		return nil, errors.New("verify passcode missing flow cookie")

	}
	authFlowParams := new(models.AuthFlowParams)
	if err := authFlowParams.ParseEncodedJson(flowCookieContents); err != nil {
		return nil, err
	}

	if err := a.structValidator.Validate(authFlowParams); err != nil {
		return nil, err
	}

	return authFlowParams, nil
}

func (a *Handler) parsePasscode(ctx *fiber.Ctx) (*models.PasscodeForm, error) {
	formValues := endpoint.ToURLValues(ctx.Context().PostArgs())
	if len(formValues) == 0 {
		return nil, errors.New("missing verify passcode request args")
	}

	passcodeForm := new(models.PasscodeForm)
	if err := a.formDecoder.Decode(passcodeForm, formValues); err != nil {
		return nil, err
	}
	if err := a.structValidator.Validate(passcodeForm); err != nil {
		return nil, err
	}

	return passcodeForm, nil
}

func renderVerifyPasscodeError(ctx *fiber.Ctx, errorMsg string) error {
	ctx.Response().Header.Set("HX-Retarget", "#verify-passcode-form .error-msg")
	return templates.Render(ctx, errroTemplates.ErrorMessage(errorMsg))
}

func renderAuthenticateInit(ctx *fiber.Ctx, errorMsg string) error {
	ctx.Response().Header.Set("HX-Retarget", "body")
	return templates.Render(ctx, pages.Authenticate(errorMsg))
}

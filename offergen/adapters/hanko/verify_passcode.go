package hanko

import (
	"fmt"
	"net/http"
	authEndpoint "offergen/endpoint/auth"
	"offergen/endpoint/models"
	"offergen/utils"

	"github.com/valyala/fasthttp"
)

func (a *Authenticator) VerifyPasscode(input *models.VerifyPasscodeInput) (*fasthttp.Cookie, error) {
	var req *fasthttp.Request
	switch input.AuthFlowParams.FlowType {
	case models.FlowTypeRegister:
		req = a.createRegisterFlowRequest(
			"verify_passcode",
			fmt.Sprintf(`{"code": "%s"}`, input.Passcode),
			input.AuthFlowParams,
		)
	case models.FlowTypeLogin:
		req = a.createLoginFlowRequest(
			"verify_passcode",
			fmt.Sprintf(`{"code": "%s"}`, input.Passcode),
			input.AuthFlowParams,
		)
	default:
		panic("flow type invalid")
	}
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.verifyPasscodeTimeout)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Info("verify passcode requst returned an error", "errMsg", err.Error())
		return nil, authEndpoint.ErrVerifyEmailRequest
	}
	switch resp.StatusCode() {
	case http.StatusOK:
	case http.StatusBadRequest:
		var verifyPasscodeResponse FlowResponse[PasscodeActions]
		utils.MustUnmarshal(resp.Body(), &verifyPasscodeResponse)
		errInvalidPasscode := models.ErrInvalidPasscode{
			Err:       "invalid passcode",
			CsrfToken: verifyPasscodeResponse.CsrfToken,
		}

		return nil, errInvalidPasscode
	default:
		logger.Info(
			"verify passcode request returned with non-200 response",
			"statusCode", resp.StatusCode(),
		)

		return nil, authEndpoint.ErrVerifyEmailRequest
	}

	authCookie := fasthttp.AcquireCookie()
	authCookie.SetKey(a.cookieName)
	if !resp.Header.Cookie(authCookie) {
		logger.Info("verify passcode request returned no cookie")

		return nil, authEndpoint.ErrVerifyEmailRequest
	}

	return authCookie, nil
}

package hanko

import (
	"errors"
	"fmt"
	"net/http"
	authEndpoint "offergen/endpoint/auth"
	"offergen/endpoint/models"
	"offergen/utils"

	"github.com/valyala/fasthttp"
)

func (a Authenticator) LoginInit(email string) (*models.AuthFlowParams, error) {
	logger.Info("starting login flow", "email", email)
	flowParams := &models.AuthFlowParams{
		Email:     email,
		FlowType:  models.FlowTypeLogin,
		CsrfToken: "",
		FlowID:    "",
	}

	if err := a.initLoginFlow(flowParams); err != nil {
		return nil, err
	}
	if err := a.registerClientCapabilitiesForLogin(flowParams); err != nil {
		return nil, err
	}
	if err := a.continueWithLoginIdentifier(email, flowParams); err != nil {
		return nil, err
	}

	return flowParams, nil
}

func (a Authenticator) initLoginFlow(flowParams *models.AuthFlowParams) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/login")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.initLoginFlowTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "init login flow",
			"errMsg", err.Error(),
		)

		return authEndpoint.ErrLoginRequest
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error(
			"response returned non-200 status code",
			"requestType", "init login flow",
			"statusCode", resp.StatusCode(),
		)

		return authEndpoint.ErrLoginRequest
	}

	var initResponse FlowResponse[InitResponseActions]
	utils.MustUnmarshal(resp.Body(), &initResponse)
	flowParams.CsrfToken = initResponse.CsrfToken
	if initResponse.CsrfToken == "" {
		return errors.New("received empty csrf token")
	}
	flowParams.MustParseFlowID(initResponse.Actions.RegisterClientCapabilities.Href)
	logger.Info("initialized login flow", "flowID", flowParams.FlowID)

	return nil
}

func (a Authenticator) registerClientCapabilitiesForLogin(flowParams *models.AuthFlowParams) error {
	req := a.createLoginFlowRequest(
		"register_client_capabilities",
		`{
					"webauthEndpointn_available": false,
					"webauthEndpointn_conditional_mediation_available": false
				}`,
		flowParams,
	)
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.registerClientCapabilitiesForLoginTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "register client capabilities for login",
			"errMsg", err.Error(),
		)

		return authEndpoint.ErrLoginRequest
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error(
			"response returned non-200 status code",
			"requestType", "register client capabilities for login",
			"statusCode", resp.StatusCode(),
		)

		return authEndpoint.ErrLoginRequest
	}

	var registerResponse FlowResponse[RegisterClientCapabilitiesActions]
	utils.MustUnmarshal(resp.Body(), &registerResponse)
	flowParams.CsrfToken = registerResponse.CsrfToken
	flowParams.MustParseFlowID(registerResponse.Actions.ContinueWithLoginIdentifier.Href)
	logger.Info("registered client capabilities", "flowID", flowParams.FlowID)

	return nil
}

func (a Authenticator) continueWithLoginIdentifier(email string, flowParams *models.AuthFlowParams) error {
	req := a.createLoginFlowRequest(
		"continue_with_login_identifier",
		fmt.Sprintf(`{"email": "%s"}`, email),
		flowParams,
	)
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.continueWithLoginIdentifierTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "continue with login identifier",
			"errMsg", err.Error(),
		)

		return authEndpoint.ErrLoginRequest
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error(
			"response returned non-200 status code",
			"requestType", "continue with login identifier",
			"statusCode", resp.StatusCode(),
		)

		return authEndpoint.ErrLoginRequest
	}

	var registerResponse FlowResponse[PasscodeActions]
	utils.MustUnmarshal(resp.Body(), &registerResponse)
	flowParams.CsrfToken = registerResponse.CsrfToken
	flowParams.MustParseFlowID(registerResponse.Actions.VerifyPasscode.Href)
	logger.Info("continue with login identifier step done", "flowID", flowParams.FlowID)

	return nil
}

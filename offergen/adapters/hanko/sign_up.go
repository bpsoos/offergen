package hanko

import (
	"fmt"
	"net/http"
	authEndpoint "offergen/endpoint/auth"
	"offergen/endpoint/models"
	"offergen/utils"

	"github.com/valyala/fasthttp"
)

func (a Authenticator) SignUp(email string) (*models.AuthFlowParams, error) {
	flowParams := &models.AuthFlowParams{
		Email:     email,
		FlowType:  models.FlowTypeRegister,
		CsrfToken: "",
		FlowID:    "",
	}
	err := a.initRegistrationFlow(flowParams)
	if err != nil {
		return nil, err
	}
	if err := a.registerClientCapabilitiesForRegister(flowParams); err != nil {
		return nil, err
	}

	if err := a.registerLoginIdentifier(email, flowParams); err != nil {
		return nil, err
	}

	return flowParams, nil
}

func (a Authenticator) initRegistrationFlow(flowParams *models.AuthFlowParams) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/registration")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.initLoginFlowTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"init registration request failed",
			"errMsg", err.Error(),
		)

		return authEndpoint.ErrSignUpRequest
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error(
			"init registration response returned non-200 status code",
			"statusCode", resp.StatusCode(),
		)

		return authEndpoint.ErrSignUpRequest
	}

	var initResponse FlowResponse[InitResponseActions]
	utils.MustUnmarshal(resp.Body(), &initResponse)
	flowParams.CsrfToken = initResponse.CsrfToken
	flowParams.MustParseFlowID(initResponse.Actions.RegisterClientCapabilities.Href)

	logger.Info("initialized registration flow", "flowID", flowParams.FlowID)

	return nil
}

func (a Authenticator) registerClientCapabilitiesForRegister(flowParams *models.AuthFlowParams) error {
	req := a.createRegisterFlowRequest(
		"register_client_capabilities",
		`{
					"webauthEndpointn_available": false,
					"webauthEndpointn_conditional_mediation_available": false
				}`,
		flowParams,
	)
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.registerClientCapabilitiesForRegisterTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "register client capabilities for register",
			"errMsg", err.Error(),
		)

		return authEndpoint.ErrSignUpRequest
	}
	if resp.StatusCode() != http.StatusOK {
		logger.Error(
			"init registration response returned non-200 status code",
			"requestType", "register client capabilities for register",
			"statusCode", resp.StatusCode(),
		)

		return authEndpoint.ErrSignUpRequest
	}

	var registerResponse FlowResponse[RegisterClientCapabilitiesActions]
	utils.MustUnmarshal(resp.Body(), &registerResponse)
	flowParams.CsrfToken = registerResponse.CsrfToken
	flowParams.MustParseFlowID(registerResponse.Actions.RegisterLoginIdentifier.Href)
	logger.Info("registered client capabilities", "flowID", flowParams.FlowID)

	return nil
}

func (a Authenticator) registerLoginIdentifier(email string, flowParams *models.AuthFlowParams) error {
	req := a.createRegisterFlowRequest(
		"register_login_identifier",
		fmt.Sprintf(`{"email": "%s"}`, email),
		flowParams,
	)
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, a.registerLoginIdentifierTimeout)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "register login identifier",
			"errMsg", err.Error(),
			"flowID", flowParams.FlowID,
		)

		return authEndpoint.ErrSignUpRequest
	}
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusConflict {
			logger.Info(
				"can't sign up, user already exists",
				"requestType", "register login identifier",
				"flowID", flowParams.FlowID,
			)
			return authEndpoint.ErrUserAlreadyExists
		}
		logger.Error(
			"response returned an unhandled non-200 status code",
			"requestType", "register login identifier",
			"statusCode", resp.StatusCode(),
			"flowID", flowParams.FlowID,
		)

		return authEndpoint.ErrSignUpRequest
	}

	var registerResponse FlowResponse[PasscodeActions]
	utils.MustUnmarshal(resp.Body(), &registerResponse)
	flowParams.CsrfToken = registerResponse.CsrfToken
	flowParams.MustParseFlowID(registerResponse.Actions.VerifyPasscode.Href)
	logger.Info("registered login identifier", "flowID", flowParams.FlowID)

	return nil
}

type FlowResponse[T any] struct {
	Name      string `json:"name"`
	Status    int    `json:"status"`
	CsrfToken string `json:"csrf_token"`
	Actions   T      `json:"actions"`
}

type InitResponseActions struct {
	RegisterClientCapabilities struct {
		Href string `json:"href"`
	} `json:"register_client_capabilities"`
}

type RegisterClientCapabilitiesActions struct {
	RegisterLoginIdentifier struct {
		Href string `json:"href"`
	} `json:"register_login_identifier"`
	ContinueWithLoginIdentifier struct {
		Href string `json:"href"`
	} `json:"continue_with_login_identifier"`
}

type PasscodeActions struct {
	Back struct {
		Href string `json:"href"`
	} `json:"back"`
	ResendPasscode struct {
		Href string `json:"href"`
	} `json:"resend_passcode"`
	VerifyPasscode struct {
		Href string `json:"href"`
	} `json:"verify_passcode"`
}

package hanko

import (
	"fmt"
	"offergen/endpoint/models"
	"time"

	"github.com/valyala/fasthttp"
)

func (a Authenticator) createRegisterFlowRequest(
	action string,
	inputData string,
	flowParams *models.AuthFlowParams,
) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/registration?action=" + action + "@" + flowParams.FlowID)
	req.SetBodyString(fmt.Sprintf(`{"input_data": %s, "csrf_token": "%s"}`, inputData, flowParams.CsrfToken))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	return req
}

func (a Authenticator) createLoginFlowRequest(
	action string,
	inputData string,
	flowParams *models.AuthFlowParams,
) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/login?action=" + action + "@" + flowParams.FlowID)
	req.SetBodyString(fmt.Sprintf(`{"input_data": %s, "csrf_token": "%s"}`, inputData, flowParams.CsrfToken))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")

	return req
}

func mustParseDuration(duration string) time.Duration {
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return parsed
}

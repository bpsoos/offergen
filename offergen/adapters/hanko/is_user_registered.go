package hanko

import (
	"fmt"
	"net/http"
	authEndpoint "offergen/endpoint/auth"
	"time"

	"github.com/valyala/fasthttp"
)

func (a *Authenticator) IsUserRegistered(email string) (bool, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/user")
	req.SetBodyRaw([]byte(fmt.Sprintf(`{"email": "%s"}`, email)))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, time.Duration(300)*time.Millisecond)

	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		logger.Error(
			"request error",
			"requestType", "get user",
			"errMsg", err.Error(),
		)

		return false, authEndpoint.ErrGetUserRequest
	}
	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNotFound {
			logger.Info("user not found", "email", email)
			return false, nil
		}
		logger.Info(
			"received non-200 status code from hanko",
			"requestType", "get user",
			"statusCode", resp.StatusCode(),
		)

		return false, authEndpoint.ErrGetUserRequest
	}

	return true, nil
}

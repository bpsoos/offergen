package hanko

import (
	"errors"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

func (a Authenticator) DeleteUser(authToken []byte) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://hanko:8000/user")
	req.Header.SetMethod(fasthttp.MethodDelete)
	req.Header.SetContentType("application/json")
	req.Header.SetCookieBytesKV([]byte(a.cookieName), authToken)
	resp := fasthttp.AcquireResponse()

	err := a.client.DoTimeout(req, resp, time.Duration(700)*time.Millisecond)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	if err != nil {
		logger.Error(
			"request failed",
			"requestType", "user delete",
			"errMsg", err.Error(),
		)

		return err
	}
	if resp.StatusCode() != http.StatusNoContent {
		logger.Error(
			"response returned non-204 status code",
			"requestType", "user delete",
			"statusCode", resp.StatusCode(),
		)

		return errors.New("non 204 status code received")
	}

	return nil
}

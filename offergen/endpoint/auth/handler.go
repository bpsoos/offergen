package auth

import (
	"errors"
	"offergen/common_deps"
	"offergen/endpoint/models"
	"offergen/logging"

	"github.com/valyala/fasthttp"
)

var logger = logging.GetLogger()

type (
	Handler struct {
		authencticator     Authenticator
		formDecoder        common_deps.FormDecoder
		structValidator    common_deps.StructValidator
		flowCookiePath     string
		flowCookieName     string
		flowCookieDomain   string
		flowCookieIsSecure bool
	}

	Authenticator interface {
		LoginInit(email string) (*models.AuthFlowParams, error)
		SignUp(email string) (*models.AuthFlowParams, error)
		VerifyPasscode(*models.VerifyPasscodeInput) (*fasthttp.Cookie, error)
		IsUserRegistered(email string) (bool, error)
	}
)

var (
	ErrLoginRequest       = errors.New("error sending login request")
	ErrSignUpRequest      = errors.New("error sending sign up request")
	ErrUserAlreadyExists  = errors.New("sign up conflict, user already exists")
	ErrVerifyEmailRequest = errors.New("error sending verify email request")
	ErrGetUserRequest     = errors.New("error sending get user request")
)

type Deps struct {
	FormDecoder     common_deps.FormDecoder
	StructValidator common_deps.StructValidator
	Authenticator   Authenticator
}

type Config struct {
	FlowCookiePath     string
	FlowCookieName     string
	FlowCookieDomain   string
	FlowCookieIsSecure bool
}

func NewHandler(config *Config, deps *Deps) *Handler {
	return &Handler{
		authencticator:     deps.Authenticator,
		formDecoder:        deps.FormDecoder,
		structValidator:    deps.StructValidator,
		flowCookiePath:     config.FlowCookiePath,
		flowCookieName:     config.FlowCookieName,
		flowCookieDomain:   config.FlowCookieDomain,
		flowCookieIsSecure: config.FlowCookieIsSecure,
	}
}

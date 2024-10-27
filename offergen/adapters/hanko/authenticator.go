package hanko

import (
	"offergen/logging"
	"time"

	"github.com/valyala/fasthttp"
)

type Authenticator struct {
	client                                       Client
	cookieName                                   string
	initRegistrationFlowTimeout                  time.Duration
	registerClientCapabilitiesForRegisterTimeout time.Duration
	registerLoginIdentifierTimeout               time.Duration
	verifyPasscodeTimeout                        time.Duration
	initLoginFlowTimeout                         time.Duration
	registerClientCapabilitiesForLoginTimeout    time.Duration
	continueWithLoginIdentifierTimeout           time.Duration
}

type Client interface {
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, timeout time.Duration) error
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}

type Deps struct {
	HTTPClient Client
}

type Config struct {
	CookieName                                   string
	InitRegistrationFlowTimeout                  string
	RegisterClientCapabilitiesForRegisterTimeout string
	RegisterLoginIdentifierTimeout               string
	VerifyPasscodeTimeout                        string
	InitLoginFlowTimeout                         string
	RegisterClientCapabilitiesForLoginTimeout    string
	ContinueWithLoginIdentifierTimeout           string
}

func NewAuthenticator(config *Config, deps *Deps) *Authenticator {
	return &Authenticator{
		client:                      deps.HTTPClient,
		cookieName:                  config.CookieName,
		initRegistrationFlowTimeout: mustParseDuration(config.InitRegistrationFlowTimeout),
		registerClientCapabilitiesForRegisterTimeout: mustParseDuration(config.RegisterClientCapabilitiesForRegisterTimeout),
		registerLoginIdentifierTimeout:               mustParseDuration(config.RegisterLoginIdentifierTimeout),
		verifyPasscodeTimeout:                        mustParseDuration(config.VerifyPasscodeTimeout),
		initLoginFlowTimeout:                         mustParseDuration(config.InitLoginFlowTimeout),
		registerClientCapabilitiesForLoginTimeout:    mustParseDuration(config.RegisterClientCapabilitiesForLoginTimeout),
		continueWithLoginIdentifierTimeout:           mustParseDuration(config.ContinueWithLoginIdentifierTimeout),
	}
}

var logger = logging.GetLogger()

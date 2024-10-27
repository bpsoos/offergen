package config

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Port                string
	HTTPClient          HTTPClient
	Auth                Auth
	DevPassword         string
	RateLimiterEnabled  bool
	PostgresURL         string
	StylesMaxAgeSeconds int
	MediaMaxAgeSeconds  int
}

type HTTPClient struct {
	ReadTimeout         string
	WriteTimeout        string
	HankoTimeouts       HankoTimeouts
	MaxIdleConnDuration string
	MaxConnDuration     string
}

type HankoTimeouts struct {
	InitRegistrationFlow                  string
	RegisterClientCapabilitiesForRegister string
	RegisterLoginIdentifier               string
	VerifyPasscode                        string
	InitLoginFlow                         string
	RegisterClientCapabilitiesForLogin    string
	ContinueWithLoginIdentifier           string
}

type Auth struct {
	URLs       AuthURLs
	CookieName string
	FlowCookie FlowCookie
}

type FlowCookie struct {
	Name     string
	Path     string
	Domain   string
	IsSecure bool
}

type AuthURLs struct {
	Base string
	Jwks string
}

const DEFAULT_PORT = "80"

func NewConfig() *Config {
	var k = koanf.New(".")

	err := k.Load(env.Provider("", "__", func(s string) string { return strings.ToLower(s) }), nil)
	if err != nil {
		panic(err)
	}
	port := k.String("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	hankoRequestTimeouts := HankoTimeouts{
		InitRegistrationFlow:                  k.MustString("http_client.hanko_timeouts.init_registration_flow"),
		RegisterClientCapabilitiesForRegister: k.MustString("http_client.hanko_timeouts.register_client_capabilities_for_register"),
		RegisterLoginIdentifier:               k.MustString("http_client.hanko_timeouts.register_login_identifier"),
		VerifyPasscode:                        k.MustString("http_client.hanko_timeouts.verify_passcode"),
		InitLoginFlow:                         k.MustString("http_client.hanko_timeouts.init_login_flow"),
		RegisterClientCapabilitiesForLogin:    k.MustString("http_client.hanko_timeouts.register_client_capabilities_for_login"),
		ContinueWithLoginIdentifier:           k.MustString("http_client.hanko_timeouts.continue_with_login_identifier"),
	}

	return &Config{
		Port:               port,
		PostgresURL:        k.MustString("postgres_url"),
		RateLimiterEnabled: k.Bool("rate_limiter_enabled"),
		HTTPClient: HTTPClient{
			ReadTimeout:         k.MustString("http_client.read_timeout"),
			WriteTimeout:        k.MustString("http_client.write_timeout"),
			MaxIdleConnDuration: k.MustString("http_client.max_idle_conn_duration"),
			MaxConnDuration:     k.MustString("http_client.max_conn_duration"),
			HankoTimeouts:       hankoRequestTimeouts,
		},
		Auth: Auth{
			URLs:       createAuthURLs(k),
			CookieName: k.String("auth.cookie_name"),
			FlowCookie: FlowCookie{
				Name:     k.String("auth.flow_cookie.name"),
				Path:     k.String("auth.flow_cookie.path"),
				Domain:   k.String("auth.flow_cookie.domain"),
				IsSecure: k.Bool("auth.flow_cookie.is_secure"),
			},
		},
		DevPassword:         k.String("dev_password"),
		StylesMaxAgeSeconds: k.Int("styles_max_age_seconds"),
		MediaMaxAgeSeconds:  k.Int("media_max_age_seconds"),
	}
}

func createAuthURLs(k *koanf.Koanf) AuthURLs {
	baseURL := k.String("auth_base_url")
	if baseURL == "" {
		panic("auth_base_url missing from config")
	}

	return AuthURLs{
		Base: baseURL,
		Jwks: baseURL + "/.well-known/jwks.json",
	}
}

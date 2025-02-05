package cmd

import (
	"context"
	"offergen/adapters/forms"
	"offergen/adapters/hanko"
	jwtAdapter "offergen/adapters/jwt"
	"offergen/adapters/validation"
	"offergen/config"
	"offergen/endpoint"
	authEndpoint "offergen/endpoint/auth"
	"offergen/endpoint/inventory"
	previewEndpoint "offergen/endpoint/preview"
	"offergen/endpoint/users"
	"offergen/persistence"
	"offergen/routing"
	"offergen/service"
	"offergen/templates"
	errorTemplates "offergen/templates/components/errors"
	inventoryTemplates "offergen/templates/components/inventory"
	pageTemplates "offergen/templates/pages"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/lestrrat-go/jwx/v2/jwk"
	_ "github.com/lib/pq"
	"github.com/valyala/fasthttp"
)

type ServeCmd struct{}

func (sc *ServeCmd) Execute() {
	config := config.NewConfig()

	client := getHTTPClient(config)
	jwksCache := getJwksCache(config)
	decoder := forms.NewDecoder()
	structValidator := validation.NewStructValidator()
	authenticator := hanko.NewAuthenticator(
		&hanko.Config{
			CookieName:                  config.Auth.CookieName,
			InitRegistrationFlowTimeout: config.HTTPClient.HankoTimeouts.InitRegistrationFlow,
			RegisterClientCapabilitiesForRegisterTimeout: config.HTTPClient.HankoTimeouts.RegisterClientCapabilitiesForRegister,
			RegisterLoginIdentifierTimeout:               config.HTTPClient.HankoTimeouts.RegisterLoginIdentifier,
			VerifyPasscodeTimeout:                        config.HTTPClient.HankoTimeouts.VerifyPasscode,
			InitLoginFlowTimeout:                         config.HTTPClient.HankoTimeouts.InitLoginFlow,
			RegisterClientCapabilitiesForLoginTimeout:    config.HTTPClient.HankoTimeouts.RegisterClientCapabilitiesForLogin,
			ContinueWithLoginIdentifierTimeout:           config.HTTPClient.HankoTimeouts.ContinueWithLoginIdentifier,
		},
		&hanko.Deps{HTTPClient: client},
	)
	tokenVerifier := jwtAdapter.NewVerifier(
		&jwtAdapter.VerifierConfig{
			JwksURL:        config.Auth.URLs.Jwks,
			AuthCookieName: config.Auth.CookieName,
		},
		&jwtAdapter.Deps{Cache: jwksCache},
	)
	inventoryTemplater := inventoryTemplates.NewInventoryTemplater()
	db := getDB(config.PostgresURL)
	inventoryManager := service.NewInventoryManager(service.InventoryManagerDeps{
		ItemPersister:      persistence.NewItemPersister(db),
		InventoryPersister: persistence.NewInventoryPersister(db),
	})

	router := routing.NewRouter(
		&routing.RouterConfig{
			StylesMaxAgeSeconds: config.StylesMaxAgeSeconds,
			MediaMaxAgeSeconds:  config.MediaMaxAgeSeconds,
		},
		&routing.RouterDeps{
			Verifier: tokenVerifier,
			RootHandler: endpoint.NewHandler(
				&endpoint.Config{AuthCookieName: config.Auth.CookieName},
				&endpoint.Deps{
					Verifier:        tokenVerifier,
					FormDecoder:     decoder,
					StructValidator: structValidator,
				},
			),
			AuthHandler: authEndpoint.NewHandler(
				&authEndpoint.Config{
					FlowCookiePath:     config.Auth.FlowCookie.Path,
					FlowCookieName:     config.Auth.FlowCookie.Name,
					FlowCookieDomain:   config.Auth.FlowCookie.Domain,
					FlowCookieIsSecure: config.Auth.FlowCookie.IsSecure,
				},
				&authEndpoint.Deps{
					Authenticator:   authenticator,
					FormDecoder:     decoder,
					StructValidator: structValidator,
				},
			),
			PreviewHandler: previewEndpoint.NewHandler(
				&previewEndpoint.Deps{
					FormDecoder:     decoder,
					StructValidator: structValidator,
				},
			),
			UserHandler: users.NewHandler(
				&users.Config{AuthCookieName: config.Auth.CookieName},
				&users.Deps{
					StructValidator:  structValidator,
					TokenVerifier:    tokenVerifier,
					InventoryManager: inventoryManager,
					UserManager: service.NewUserManager(
						&service.UserManagerDeps{
							Persister:     persistence.NewUserPersister(db),
							Authenticator: authenticator,
						},
					),
				},
			),
			InventoryHandler: inventory.NewHandler(
				&inventory.Deps{
					FormDecoder:      decoder,
					StructValidator:  structValidator,
					TokenVerifier:    tokenVerifier,
					InventoryManager: inventoryManager,
					Renderer:         templates.NewRenderer(),
					ErrorTemplater:   errorTemplates.NewErrorTemplater(),
					PageTemplater: pageTemplates.NewPageTemplater(
						strconv.FormatInt(time.Now().Unix(), 10),
						&pageTemplates.Deps{
							InventoryTemplater: inventoryTemplater,
						},
					),
					InventoryTemplater: inventoryTemplater,
				},
			),
		})

	app := fiber.New(
		fiber.Config{
			ReadBufferSize: 16384,
			ProxyHeader:    "X-Forwarded-For",
		},
	)
	app.Use(getLoggerMiddleware())

	if config.RateLimiterEnabled {
		app.Use(getRateLimiterMiddleware())
	}

	router.AddRoutes(app)

	logger.Info("Starting to listen")
	if err := app.Listen(":" + config.Port); err != nil {
		panic(err)
	}
}

func getHTTPClient(config *config.Config) *fasthttp.Client {
	readTimeout, err := time.ParseDuration(config.HTTPClient.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := time.ParseDuration(config.HTTPClient.WriteTimeout)
	if err != nil {
		panic(err)
	}
	maxIdleConnDuration, err := time.ParseDuration(config.HTTPClient.MaxIdleConnDuration)
	if err != nil {
		panic(err)
	}
	maxConnDuration, err := time.ParseDuration(config.HTTPClient.MaxConnDuration)
	if err != nil {
		panic(err)
	}

	return &fasthttp.Client{
		MaxConnDuration:     maxConnDuration,
		ReadTimeout:         readTimeout,
		WriteTimeout:        writeTimeout,
		MaxIdleConnDuration: maxIdleConnDuration,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}

func getJwksCache(config *config.Config) *jwk.Cache {
	jwksCache := jwk.NewCache(context.Background())
	err := jwksCache.Register(config.Auth.URLs.Jwks)
	if err != nil {
		panic(err)
	}
	_, err = jwksCache.Refresh(context.Background(), config.Auth.URLs.Jwks)
	if err != nil {
		panic(err)
	}
	return jwksCache
}

func getDB(url string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		panic(err)
	}

	return db
}

func getLoggerMiddleware() fiber.Handler {
	return fiberLogger.New(
		fiberLogger.Config{
			Format:        `{"time": "${time}", "level": "INFO", "msg": {"accessLog": {"status": "${status}", "latency": "${latency}", "ip": "${ip}", "forwardedFor": "${reqHeader:X-Forwarded-For}", "userAgent": "${reqHeader:User-Agent}", "method": "${method}", "path": "${path}", "error": "${error}"}}}` + "\n",
			TimeFormat:    time.RFC3339Nano,
			TimeZone:      "UTC",
			DisableColors: true,
		},
		fiberLogger.ConfigDefault,
	)
}

func getRateLimiterMiddleware() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn("rate limit reached", "ip", c.IP())
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.FixedWindow{},
	})
}

func getBasicAuthMiddleware(config *config.Config) fiber.Handler {
	return basicauth.New(
		basicauth.Config{
			Users: map[string]string{
				"admin": config.DevPassword,
			},
		})
}

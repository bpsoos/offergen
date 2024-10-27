package routing

import (
	"io/fs"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	stylesMaxAgeSeconds int
	mediaMaxAgeSeconds  int
	verifier            TokenVerifier
	rootHandler         RootHandler
	authHandler         AuthHandler
	previewHandler      PreviewHandler
	userHandler         UserHandler
	inventoryHandler    InventoryHandler
}

type RouterDeps struct {
	Verifier         TokenVerifier
	RootHandler      RootHandler
	AuthHandler      AuthHandler
	PreviewHandler   PreviewHandler
	UserHandler      UserHandler
	InventoryHandler InventoryHandler
}

type TokenVerifier interface {
	VerifyUser(ctx *fiber.Ctx) error
}

type RootHandler interface {
	Index(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	Health(ctx *fiber.Ctx) error
	Authenticate(ctx *fiber.Ctx) error
	Styles() fs.FS
	Media() fs.FS
}

type AuthHandler interface {
	SignUp(ctx *fiber.Ctx) error
	VerifyPasscode(ctx *fiber.Ctx) error
	Init(ctx *fiber.Ctx) error
}

type PreviewHandler interface {
	Add(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Generate(ctx *fiber.Ctx) error
}

type InventoryHandler interface {
	Create(ctx *fiber.Ctx) error
	CreatePage(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	Editor(ctx *fiber.Ctx) error
	Items(ctx *fiber.Ctx) error
	ItemPages(ctx *fiber.Ctx) error
	ItemsPage(ctx *fiber.Ctx) error
	SettingsPage(ctx *fiber.Ctx) error
	UpdateInventory(ctx *fiber.Ctx) error
}

type UserHandler interface {
	Create(ctx *fiber.Ctx) error
	Profile(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type RouterConfig struct {
	StylesMaxAgeSeconds int
	MediaMaxAgeSeconds  int
}

func NewRouter(config *RouterConfig, deps *RouterDeps) *Router {
	return &Router{
		stylesMaxAgeSeconds: config.StylesMaxAgeSeconds,
		mediaMaxAgeSeconds:  config.MediaMaxAgeSeconds,
		rootHandler:         deps.RootHandler,
		verifier:            deps.Verifier,
		authHandler:         deps.AuthHandler,
		previewHandler:      deps.PreviewHandler,
		userHandler:         deps.UserHandler,
		inventoryHandler:    deps.InventoryHandler,
	}
}

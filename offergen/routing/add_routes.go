package routing

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func (r *Router) AddRoutes(app *fiber.App) {
	app.Get("/health", r.rootHandler.Health)
	event := app.Group("/event")
	event.Post("/users/create", r.userHandler.Create)

	app.Get("/", r.rootHandler.Index)

	app.Use("/media", filesystem.New(filesystem.Config{
		Root:               http.FS(r.rootHandler.Media()),
		MaxAge:             r.mediaMaxAgeSeconds,
		ContentTypeCharset: "image/svg+xml",
	}))
	app.Use("/styles", filesystem.New(filesystem.Config{
		Root:               http.FS(r.rootHandler.Styles()),
		MaxAge:             r.stylesMaxAgeSeconds,
		ContentTypeCharset: "text/css",
	}))

	app.Get("/authenticate", r.rootHandler.Authenticate)
	app.Get("/logout", r.rootHandler.Logout)

	auth := app.Group("/auth")
	auth.Post("/init", r.authHandler.Init)
	auth.Post("/signup", r.authHandler.SignUp)
	auth.Post("/verify-passcode", r.authHandler.VerifyPasscode)

	preview := app.Group("/preview")
	preview.Post("/items", r.previewHandler.Add)
	preview.Post("/generate", r.previewHandler.Generate)
	preview.Delete("/items/:id", r.previewHandler.Delete)

	user := app.Group("/user")
	user.Use(r.verifier.VerifyUser)
	user.Get("", r.userHandler.Profile)
	user.Delete("", r.userHandler.Delete)

	inventory := app.Group("/inventory")
	inventory.Use(r.verifier.VerifyUser)

	inventory.Get("", r.inventoryHandler.Editor)
	inventory.Get("/categories", r.inventoryHandler.BatchGetCategory)
	inventory.Post("/categories", r.inventoryHandler.CreateCategory)
	inventory.Get("/categories/init", r.inventoryHandler.CreateCategoryInit)
	inventory.Get("/settings-page", r.inventoryHandler.SettingsPage)
	inventory.Post("/update", r.inventoryHandler.UpdateInventory)

	inventory.Get("/create-page", r.inventoryHandler.CreatePage)
	inventory.Get("/items", r.inventoryHandler.Items)
	inventory.Get("/item-pages", r.inventoryHandler.ItemPages)
	inventory.Post("/item", r.inventoryHandler.CreateItem)
	inventory.Delete("/item/:id", r.inventoryHandler.Delete)

	offerings := app.Group("/o")
	offerings.Get(":owner_id", r.offeringHandler.GetOffering)
}

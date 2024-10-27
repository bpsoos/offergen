package templates

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

func Render(ctx *fiber.Ctx, component templ.Component) error {
	ctx.Type("html", "utf-8")
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

type Renderer struct{}

func (r *Renderer) Render(ctx *fiber.Ctx, component templ.Component) error {
	ctx.Type("html", "utf-8")
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

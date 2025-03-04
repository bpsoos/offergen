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

func NewRenderer() *Renderer {
	return &Renderer{}
}

func (r *Renderer) Render(
	ctx *fiber.Ctx,
	component templ.Component,
) error {
	ctx.Type("html", "utf-8")
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func (r *Renderer) RenderWithChildren(
	ctx *fiber.Ctx,
	children templ.Component,
	component templ.Component,
) error {
	ctx.Type("html", "utf-8")
	return component.Render(
		templ.WithChildren(ctx.Context(), children),
		ctx.Response().BodyWriter(),
	)
}

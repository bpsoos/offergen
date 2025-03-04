package handlermachinery

import (
	"context"
	"io"

	"github.com/gofiber/fiber/v2"
)

func ToTemplaterContext(ctx *fiber.Ctx) (context.Context, io.Writer) {
	ctx.Type("html", "utf-8")
	return ctx.Context(), ctx.Response().BodyWriter()
}

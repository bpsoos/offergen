package handlermachinery

import (
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func ParseURLValues(ctx *fiber.Ctx) url.Values {
	args := ctx.Request().PostArgs()
	if args == nil || args.Len() == 0 {
		return nil
	}

	values := make(url.Values)
	args.VisitAll(func(key, value []byte) {
		values.Add(string(key), string(value))
	})

	return values
}

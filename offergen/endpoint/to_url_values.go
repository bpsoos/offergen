package endpoint

import (
	"net/url"

	"github.com/valyala/fasthttp"
)

func ToURLValues(args *fasthttp.Args) url.Values {
	if args == nil || args.Len() == 0 {
		return nil
	}

	values := make(url.Values)
	args.VisitAll(func(key, value []byte) {
		values.Add(string(key), string(value))
	})

	return values
}

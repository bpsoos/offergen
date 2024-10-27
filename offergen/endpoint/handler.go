package endpoint

import (
	"context"
	"offergen/common_deps"
)

type Handler struct {
	verifier        TokenVerifier
	structValidator common_deps.StructValidator
	formDecoder     common_deps.FormDecoder
	authCookieName  string
}

type Config struct {
	AuthCookieName string
}

type Deps struct {
	Verifier        TokenVerifier
	FormDecoder     common_deps.FormDecoder
	StructValidator common_deps.StructValidator
}

type TokenVerifier interface {
	IsValidToken(ctx context.Context, token []byte) bool
}

func NewHandler(config *Config, deps *Deps) *Handler {
	return &Handler{
		verifier:        deps.Verifier,
		structValidator: deps.StructValidator,
		formDecoder:     deps.FormDecoder,
		authCookieName:  config.AuthCookieName,
	}
}

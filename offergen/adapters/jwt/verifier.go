package jwt

import (
	"context"
	"crypto/rsa"
	"offergen/logging"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type Verifier struct {
	jwksURL        string
	authCookieName string
	cache          *jwk.Cache
}

type JwksCache interface {
	Get(ctx context.Context, url string) (KeySet, error)
	Refresh(ctx context.Context, url string) (KeySet, error)
}

type KeySet interface {
	LookupKeyID(string) (Key, bool)
}

type Key interface {
	Raw(*rsa.PublicKey) error
}

type Parser interface {
	Parse(s []byte) (Token, error)
}

type Token interface {
	Subject() string
}

type Deps struct {
	Cache *jwk.Cache
}

type VerifierConfig struct {
	JwksURL        string
	AuthCookieName string
}

func NewVerifier(config *VerifierConfig, deps *Deps) *Verifier {
	return &Verifier{
		jwksURL:        config.JwksURL,
		cache:          deps.Cache,
		authCookieName: config.AuthCookieName,
	}
}

var logger = logging.GetLogger()

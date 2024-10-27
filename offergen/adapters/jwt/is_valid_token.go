package jwt

import (
	"context"
	"errors"

	jwxJwt "github.com/lestrrat-go/jwx/v2/jwt"
)

func (v Verifier) IsValidToken(ctx context.Context, rawToken []byte) bool {
	_, err := v.parseToken(ctx, rawToken)

	return err == nil
}

func (v *Verifier) parseToken(ctx context.Context, rawToken []byte) (jwxJwt.Token, error) {
	if len(rawToken) == 0 {
		return nil, errors.New("token lenght was 0")
	}

	keyset, err := v.cache.Get(ctx, v.jwksURL)
	if err != nil {
		logger.Error(
			"could not get keyset",
			"url", v.jwksURL,
			"errMsg", err.Error(),
		)
		panic("could not get")
	}

	logger.Info("validating token")
	token, err := jwxJwt.Parse(rawToken, jwxJwt.WithKeySet(keyset))
	if err == nil {
		return token, nil
	}

	logger.Info("could not validate token, trying to refresh cache")
	keyset, err = v.cache.Refresh(ctx, v.jwksURL)
	if err != nil {
		logger.Error(
			"could not get keyset",
			"url", v.jwksURL,
			"errMsg", err.Error(),
		)
		panic("could not get")
	}
	token, err = jwxJwt.Parse(rawToken, jwxJwt.WithKeySet(keyset))
	if err != nil {
		logger.Info("invalid token")

		return nil, errors.New("invalid token")
	}
	logger.Info("validated token", "subject", token.Subject())

	return token, nil
}

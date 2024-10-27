package jwt

import (
	"context"
	"errors"
)

func (v *Verifier) GetWebhookClaims(rawToken string) (data map[string]interface{}, evt string, err error) {
	token, err := v.parseToken(context.TODO(), []byte(rawToken))
	if err != nil {
		return nil, "", err
	}

	rawData, dataFound := token.Get("data")
	if !dataFound {
		return nil, "", errors.New("missing data key")
	}
	rawEvt, evtFound := token.Get("evt")
	if !evtFound {
		return nil, "", errors.New("missing evt key")
	}

	return rawData.(map[string]interface{}), rawEvt.(string), nil
}

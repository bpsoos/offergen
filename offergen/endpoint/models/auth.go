package models

import (
	"encoding/base64"
	"net/url"
	"offergen/utils"
	"strings"
)

type EmailForm struct {
	Email string `form:"Email" validate:"required,email"`
}

type PasscodeForm struct {
	Passcode string `form:"Passcode" validate:"required,numeric,len=6"`
}

type VerifyPasscodeInput struct {
	*AuthFlowParams
	*PasscodeForm
}

type FlowType string

const (
	FlowTypeLogin    FlowType = "login"
	FlowTypeRegister FlowType = "register"
)

type AuthFlowParams struct {
	CsrfToken string   `json:"csrf_token" validate:"required,alphanum,len=32"`
	FlowID    string   `json:"flow_id" validate:"required,uuid"`
	Email     string   `json:"email" validate:"required,email"`
	FlowType  FlowType `json:"flow_type" validate:"required,oneof=login register"`
}

func (afp *AuthFlowParams) MustParseFlowID(href string) {
	parsedURL, _ := url.Parse(href)
	splitAction := strings.Split(parsedURL.Query().Get("action"), "@")
	if len(splitAction) != 2 {
		panic("could not parse id from href")
	}
	afp.FlowID = splitAction[1]
}

func (afp *AuthFlowParams) ToEncodedJson() []byte {
	value := utils.MustMarshal(afp)
	encodedValue := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
	base64.StdEncoding.Encode(encodedValue, value)

	return encodedValue
}

func (afp *AuthFlowParams) ParseEncodedJson(encodedJsonParams []byte) error {
	decodedParams := make([]byte, base64.StdEncoding.DecodedLen(len(encodedJsonParams)))
	written, err := base64.StdEncoding.Decode(decodedParams, encodedJsonParams)
	if err != nil {

		return err
	}
	utils.MustUnmarshal(decodedParams[:written], afp)

	return nil
}

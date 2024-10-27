package models_test

import (
	"offergen/endpoint/models"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("flow response", func() {
	Context("parse flow id called", func() {
		It("should parse flow id from href", func() {
			flowID := uuid.NewString()
			params := models.AuthFlowParams{}
			params.MustParseFlowID("/registration?action=register_client_capabilities%40" + flowID)

			Expect(params.FlowID).To(Equal(flowID))
		})
	})
	Context("params need to be encoded and decoded", func() {
		It("should return the original values after encoding and decoding", func() {
			params := models.AuthFlowParams{
				CsrfToken: "dBVYDZlkw5NaTkv93e3KxHpZP0EyxQiN",
				FlowID:    "bac8ba09-b1de-4819-bf2c-c6607e52aabe",
				FlowType:  models.FlowTypeRegister,
				Email:     "dummy@email.com",
			}

			var decodedParams models.AuthFlowParams
			err := decodedParams.ParseEncodedJson(params.ToEncodedJson())
			Expect(err).To(Not(HaveOccurred()))

			Expect(decodedParams.CsrfToken).To(Equal("dBVYDZlkw5NaTkv93e3KxHpZP0EyxQiN"))
			Expect(decodedParams.FlowID).To(Equal("bac8ba09-b1de-4819-bf2c-c6607e52aabe"))
			Expect(decodedParams.FlowType).To(Equal(models.FlowTypeRegister))
			Expect(decodedParams.Email).To(Equal("dummy@email.com"))
		})

		It("should return the original values of various inputs after encoding and decoding", func() {
			params := models.AuthFlowParams{
				CsrfToken: "ls8UnWRpekoZXW7Dh8RIsz4U1yeFxD7Y",
				FlowID:    uuid.NewString(),
				FlowType:  models.FlowTypeLogin,
				Email:     "other_dummy@email.com",
			}
			encodedJsonParams := params.ToEncodedJson()

			var decodedParams models.AuthFlowParams
			err := decodedParams.ParseEncodedJson(encodedJsonParams)
			Expect(err).To(Not(HaveOccurred()))

			Expect(decodedParams.CsrfToken).To(Equal(params.CsrfToken))
			Expect(decodedParams.FlowID).To(Equal(params.FlowID))
			Expect(decodedParams.FlowType).To(Equal(models.FlowTypeLogin))
			Expect(decodedParams.Email).To(Equal("other_dummy@email.com"))
		})
	})
},
)

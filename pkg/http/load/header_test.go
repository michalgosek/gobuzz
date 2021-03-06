package load_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/gobuzz/pkg/http/load"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// testContent is an internal aggregate for creating tableTest slice.
type testContent struct {
	payloadContent   *strings.Reader
	validationResult PayloadValidationError
	contentType      string
}

func fakeHandler(w http.ResponseWriter, r *http.Request) PayloadValidationError {
	var checkStruct JSONPostBody
	result := PostPayloadCheck(w, r, &checkStruct)
	return result
}

var _ = Describe("When calling PostPayloadCheck", func() {

	var data []testContent
	BeforeEach(func() {
		data = []testContent{
			{
				strings.NewReader(`{"url": "https://httpbin.org/range/15","interval":60}`),
				PayloadValidationError{Status: http.StatusAccepted, Msg: fmt.Sprintln("Payload check validation was succed.")},
				"application/json",
			},
			{
				strings.NewReader(`{"url": "https://httpbin.org/delay/15","interval":110}`),
				PayloadValidationError{Status: http.StatusAccepted, Msg: fmt.Sprintln("Payload check validation was succed.")},
				"application/json",
			},
		}
	})

	Context("When JSON payload is valid.", func() {
		It("Should return http.StatusOK, msg: validation ok.", func() {
			for _, el := range data {
				r := httptest.NewRequest(http.MethodPost, "/", el.payloadContent)
				w := httptest.NewRecorder()
				result := fakeHandler(w, r)
				Expect(result.Status).To(Equal(el.validationResult.Status))
				Expect(result.Msg).To(Equal(el.validationResult.Msg))
			}
		})
	})

	Context("When JSON payload is invalid.", func() {
		BeforeEach(func() {
			data = []testContent{
				{
					strings.NewReader(`{"url": "https://httpbin.org/range/40"}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Missing interval field in JSON payload.")},
					"application/json",
				},
				{
					strings.NewReader(`{"url": "https://httpbin.org/range/40", "interval":10}Woops!'`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body must cotain only single JSON object.")},
					"application/json",
				},
				{
					strings.NewReader(`{"url": 1234, "interval": ""}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln(`Request body contains an invalid value for the "url" field (at position 12).`)},
					"application/json",
				},
				{
					strings.NewReader(`{"url": "http://httpbin.org/range/15","interval":"abc"}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln(`Request body contains an invalid value for the "interval" field (at position 54).`)},
					"application/json",
				},
				{
					strings.NewReader(``),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body must not be empty.")},
					"application/json",
				},
				{
					strings.NewReader(`{}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Missing url and interval fields in JSON payload.")},
					"application/json",
				},
				{
					strings.NewReader(`{"interval":10}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Missing url field in JSON payload.")},
					"application/json",
				},
				{
					strings.NewReader(`{interval:""}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body contains badly-formed JSON (at position 2).")},
					"application/json",
				},
				{
					strings.NewReader(`{interval:}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body contains badly-formed JSON (at position 2).")},
					"application/json",
				},
				{
					strings.NewReader(`{"url":"abc"}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Missing interval field in JSON payload.")},
					"application/json",
				},
				{
					strings.NewReader(`{url:""}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body contains badly-formed JSON (at position 2).")},
					"application/json",
				},
				{
					strings.NewReader(`{url:}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln("Request body contains badly-formed JSON (at position 2).")},
					"application/json",
				},
				{
					strings.NewReader(`{"url": "https://httpbin.org/range/15","interval":60, "zonk": "Hello!"}`),
					PayloadValidationError{Status: http.StatusBadRequest, Msg: fmt.Sprintln(`Request body contains unknown field "zonk".`)},
					"application/json",
				},
				{ // Invalid Content-Types:
					strings.NewReader(`{"url": "https://httpbin.org/range/15","interval":60, "zonk": "Hello!"}`),
					PayloadValidationError{Status: http.StatusUnsupportedMediaType, Msg: fmt.Sprintln("Invalid or lack of Content-Type.")},
					"text/html",
				},
				{
					strings.NewReader(`{"url": "https://httpbin.org/range/15","interval":60, "zonk": "Hello!"}`),
					PayloadValidationError{Status: http.StatusUnsupportedMediaType, Msg: fmt.Sprintln("Invalid or lack of Content-Type.")},
					"application/x-www-form-urlencoded",
				},
			}
		})

		It("Should return correct httpStatus and validation msg.", func() {
			for _, el := range data {
				r := httptest.NewRequest(http.MethodPost, "/", el.payloadContent)
				r.Header.Set("Content-Type", el.contentType)
				w := httptest.NewRecorder()
				result := fakeHandler(w, r)
				Expect(result.Status).To(Equal(el.validationResult.Status))
				Expect(result.Error()).To(Equal(el.validationResult.Msg))
			}
		})
	})
})

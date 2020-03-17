package load

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

// JSONPostBody represents an HTTP fetch request recived by server
// from the client. Field describes value of each key.
// Used for decoding request body operation.
type JSONPostBody struct {
	URL      *string `json:"url"`
	Interval *int    `json:"interval"`
}

// Validate reports wether sending JSON payload has valid structure
// containing url and interval fields. If so, method returns nil which
// indicates confirmation. Otherwise returns http status code and suggestion
// text back to the client if payload has missing fields.
func (j *JSONPostBody) Validate() PayloadValidationError {
	switch {
	case j.URL == nil && j.Interval == nil:
		txt := fmt.Sprintln("Missing url and interval fields in JSON payload.")
		return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}
	case j.URL == nil:
		txt := fmt.Sprintln("Missing url field in JSON payload.")
		return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}
	case j.Interval == nil:
		txt := fmt.Sprintln("Missing interval field in JSON payload.")
		return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}
	}
	txt := fmt.Sprintln("url, interval fields validation was succed.")
	return PayloadValidationError{Status: http.StatusAccepted, Msg: txt}
}

// DecodeHandleReport provides a contextual information to the client
// about badly formated JSON sent payload if any error occurs.
// If so, method returns http status code and error msg describing
// different case. Otherwise returns nil to indicate decoding was succeed.
//
// TO-DO: think about shorter msg returned back to Client.
func (j *JSONPostBody) DecodeHandleReport(err error) PayloadValidationError {
	var (
		syntaxError        *json.SyntaxError
		unmarshalTypeError *json.UnmarshalTypeError
	)

	if err != nil {
		switch {
		// Go 1.13 - https://blog.golang.org/go1.13-errors
		case errors.As(err, &syntaxError):
			txt := fmt.Sprintf("Request body contains badly-formed JSON (at position %d).\n", syntaxError.Offset)
			return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}

		case errors.Is(err, io.ErrUnexpectedEOF):
			txt := fmt.Sprintln("Request body contains badly-formed JSON.")
			return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}

		case errors.As(err, &unmarshalTypeError):
			txt := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d).\n", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			txt := fmt.Sprintf("Request body contains unknown field %s.\n", fieldName)
			return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}

		case errors.Is(err, io.EOF):
			txt := fmt.Sprintln("Request body must not be empty.")
			return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}

		case err.Error() == "http: request body too large":
			txt := fmt.Sprintln("Request body must not be larger than 1MB.")
			return PayloadValidationError{Status: http.StatusRequestEntityTooLarge, Msg: txt}

		default:
			log.Println(err.Error())
			return PayloadValidationError{Status: http.StatusInternalServerError, Msg: http.StatusText(http.StatusInternalServerError)}
		}
	}

	txt := fmt.Sprintln("Decoding validation was succed.")
	return PayloadValidationError{Status: http.StatusAccepted, Msg: txt}
}

// PayloadValidationError represetns response body sending to client
// when any error occurs during PostPayloadCheck run.
type PayloadValidationError struct {
	Status int
	Msg    string
}

// Error returns specific error message to the Client.
func (p *PayloadValidationError) Error() string {
	return p.Msg
}

// PostPayloadCheck returns error response if sending payload is larger
// than 1MB and if the JSON post body is not formatted correctly.
// Decoded requst body data is saved in content arg.
func PostPayloadCheck(w http.ResponseWriter, r *http.Request, content *JSONPostBody) PayloadValidationError {
	if r.Header.Get("Content-Type") != "" {
		if val, _ := header.ParseValueAndParams(r.Header, "Content-Type"); val != "application/json" {
			return PayloadValidationError{Status: http.StatusUnsupportedMediaType, Msg: "Invalid or lack of Content-Type"}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // Limit to 1MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // Unwanted fields check
	err := dec.Decode(&content)

	// Decode error handling
	if err != nil {
		return content.DecodeHandleReport(err)
	}

	// Extraneous json data in request body
	if dec.More() {
		txt := fmt.Sprintln("Request body must cotain only single JSON object.")
		return PayloadValidationError{Status: http.StatusBadRequest, Msg: txt}
	}

	// Content validation
	validState := content.Validate()
	if validState.Status != http.StatusAccepted {
		return validState
	}

	txt := fmt.Sprintln("Payload check validation was succed.")
	return PayloadValidationError{Status: http.StatusAccepted, Msg: txt}
}

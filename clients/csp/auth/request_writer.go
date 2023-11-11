package auth

import (
	"context"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// RequestWriter contains all the parameters required to retrieve a new access token.
type RequestWriter struct {
	Body *RequestBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WriteToRequest writes these params to a request
func (o *RequestWriter) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	return nil
}

// RequestBody contains all the parameters required to retrieve a new auth token.
type RequestBody struct {
	// Access token: https://docs.wavefront.com/using_wavefront_api.html#make-api-calls-by-using-a-user-account
	// Example: whatM9vIgu4sa2bytesTRows6NZR8QxAL3vQJ6QcGLaTRKvT0jLDogfishFJRA32
	// Required: true
	ApiToken *string `json:"api_token"`

	// TODO: probably need to support an mfa authentication technique.
	// What is a standard way to get user input when using terraform?
	// Seems like we need to start caching the access token so they don't have to repeat mfa.
	// I am not building with this in mind currently. Trying to stay simple and not over build.
	// Required true if using an MFA device.
	// Passcode *string `json:"passcode,omitempty"`
}

// Validate validates this csp login specification
func (m *RequestBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRefreshToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RequestBody) validateRefreshToken(formats strfmt.Registry) error {

	if err := validate.Required("api_token", "body", m.ApiToken); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this csp login specification based on context it is used
func (m *RequestBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RequestBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RequestBody) UnmarshalBinary(b []byte) error {
	var res RequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

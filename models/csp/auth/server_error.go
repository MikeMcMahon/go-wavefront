package auth

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ServerError Entity used when the server encounters an unexpected condition,
// preventing it from fulfilling the request. (5xx)
type ServerError struct {
	Uri      string `json:"uri"`
	CspError struct {
		ErrorCode         int      `json:"errorCode"`
		Metadata          struct{} `json:"metadata"`
		NumericModuleCode int      `json:"numericModuleCode"`
		ModuleCode        string   `json:"moduleCode"`
		ErrorMessageCode  string   `json:"errorMessageCode"`
	} `json:"cspError"`
	Causes []string `json:"causes"`
	// ServerError message
	// Example: unknown error occurred.
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestId"`
	// status code
	StatusCode int `json:"statusCode,omitempty"`
}

// Validate validates this service error response
func (m *ServerError) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this error based on the context it is used in
func (m *ServerError) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ServerError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServerError) UnmarshalBinary(b []byte) error {
	var res ServerError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

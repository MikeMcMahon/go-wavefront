package auth

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ClientError Entity used when the server cannot or will not process a request,
// due to something that is perceived to be a client error. (4xx)
type ClientError struct {
	Uri      string `json:"uri"`
	CspError struct {
		ErrorCode         int      `json:"errorCode"`
		Metadata          struct{} `json:"metadata"`
		NumericModuleCode int      `json:"numericModuleCode"`
		ModuleCode        string   `json:"moduleCode"`
		ErrorMessageCode  string   `json:"errorMessageCode"`
	} `json:"cspError"`
	Causes []string `json:"causes"`
	// ClientError message
	// Example: Failed to validate credentials.
	Message   string `json:"message,omitempty"`
	RequestId string `json:"requestId"`
	// status code
	StatusCode int `json:"statusCode,omitempty"`
}

// Validate validates this error
func (m *ClientError) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this error based on the context it is used in
func (m *ClientError) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ClientError) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ClientError) UnmarshalBinary(b []byte) error {
	var res ClientError
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

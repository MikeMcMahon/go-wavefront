package auth

import (
	"fmt"
	auth "github.com/WavefrontHQ/go-wavefront-management-api/v2/models/csp/auth"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// ResponseReader handles auth related server responses.
type ResponseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the appropriate response object.
func (o *ResponseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := &OK200{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil // Only case with a non error response.
	case 400:
		result := &InvalidRequest400{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := &NotFound404{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := &Conflict409{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := &TooManyRequests429{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := &Unexpected500{}
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("server responded with unknown status code", response, response.Code())
	}
}

type OK200 struct{ Payload *auth.Success }

func (o *OK200) IsSuccess() bool      { return true }
func (o *OK200) IsRedirect() bool     { return false }
func (o *OK200) IsClientError() bool  { return false }
func (o *OK200) IsServerError() bool  { return false }
func (o *OK200) IsCode(code int) bool { return code == 200 }

func (o *OK200) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] OK200  %+v", 200, o.Payload)
}

func (o *OK200) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] OK200  %+v", 200, o.Payload)
}

func (o *OK200) GetPayload() *auth.Success {
	return o.Payload
}

func (o *OK200) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.Success)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type InvalidRequest400 struct{ Payload *auth.ClientError }

func (o *InvalidRequest400) IsSuccess() bool      { return false }
func (o *InvalidRequest400) IsRedirect() bool     { return false }
func (o *InvalidRequest400) IsClientError() bool  { return true }
func (o *InvalidRequest400) IsServerError() bool  { return false }
func (o *InvalidRequest400) IsCode(code int) bool { return code == 400 }

func (o *InvalidRequest400) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] InvalidRequest400  %+v", 400, o.Payload)
}

func (o *InvalidRequest400) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] InvalidRequest400  %+v", 400, o.Payload)
}

func (o *InvalidRequest400) GetPayload() *auth.ClientError {
	return o.Payload
}

func (o *InvalidRequest400) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.ClientError)

	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type NotFound404 struct{ Payload *auth.ClientError }

func (o *NotFound404) IsSuccess() bool      { return false }
func (o *NotFound404) IsRedirect() bool     { return false }
func (o *NotFound404) IsClientError() bool  { return true }
func (o *NotFound404) IsServerError() bool  { return false }
func (o *NotFound404) IsCode(code int) bool { return code == 404 }

func (o *NotFound404) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] NotFound404  %+v", 404, o.Payload)
}

func (o *NotFound404) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] NotFound404  %+v", 404, o.Payload)
}

func (o *NotFound404) GetPayload() *auth.ClientError {
	return o.Payload
}

func (o *NotFound404) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.ClientError)

	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type Conflict409 struct{ Payload *auth.ClientError }

func (o *Conflict409) IsSuccess() bool      { return false }
func (o *Conflict409) IsRedirect() bool     { return false }
func (o *Conflict409) IsClientError() bool  { return true }
func (o *Conflict409) IsServerError() bool  { return false }
func (o *Conflict409) IsCode(code int) bool { return code == 409 }

func (o *Conflict409) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] Conflict409  %+v", 409, o.Payload)
}

func (o *Conflict409) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] Conflict409  %+v", 409, o.Payload)
}

func (o *Conflict409) GetPayload() *auth.ClientError {
	return o.Payload
}

func (o *Conflict409) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.ClientError)

	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type TooManyRequests429 struct{ Payload *auth.ClientError }

func (o *TooManyRequests429) IsSuccess() bool      { return false }
func (o *TooManyRequests429) IsRedirect() bool     { return false }
func (o *TooManyRequests429) IsClientError() bool  { return true }
func (o *TooManyRequests429) IsServerError() bool  { return false }
func (o *TooManyRequests429) IsCode(code int) bool { return code == 429 }

func (o *TooManyRequests429) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] TooManyRequests429  %+v", 429, o.Payload)
}

func (o *TooManyRequests429) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] TooManyRequests429  %+v", 429, o.Payload)
}

func (o *TooManyRequests429) GetPayload() *auth.ClientError {
	return o.Payload
}

func (o *TooManyRequests429) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.ClientError)

	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type Unexpected500 struct{ Payload *auth.ServerError }

func (o *Unexpected500) IsSuccess() bool      { return false }
func (o *Unexpected500) IsRedirect() bool     { return false }
func (o *Unexpected500) IsClientError() bool  { return false }
func (o *Unexpected500) IsServerError() bool  { return true }
func (o *Unexpected500) IsCode(code int) bool { return code == 500 }

func (o *Unexpected500) Error() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] Unexpected500  %+v", 500, o.Payload)
}

func (o *Unexpected500) String() string {
	return fmt.Sprintf("[POST /csp/gateway/am/api/auth/api-tokens/authorize][%d] Unexpected500  %+v", 500, o.Payload)
}

func (o *Unexpected500) GetPayload() *auth.ServerError {
	return o.Payload
}

func (o *Unexpected500) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(auth.ServerError)

	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

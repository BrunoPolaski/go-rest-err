package rest_err

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RestErr struct {
	Message   string    `json:"message" example:"invalid request parameters"` // Human readable message
	Err       string    `json:"error" example:"bad request"`
	Code      int       `json:"code" example:"400"` // HTTP status code
	Causes    []Causes  `json:"causes,omitempty"`   // Detailed error causes, most common for json field validation errors
	Timestamp time.Time `json:"timestamp"`          // When the error occurred
	Wrapped   error     `json:"-"`                  // Underlying error (not exposed in JSON)
}

type Causes struct {
	Field   string `json:"field" example:"email"`                   // Field or parameter that caused the error
	Message string `json:"message" example:"invalid email address"` // Description of the cause
}

func (r *RestErr) Error() string {
	if r.Wrapped != nil {
		return fmt.Sprintf("%s: %v", r.Message, r.Wrapped)
	}
	return r.Message
}

// Unwrap returns the underlying error for error chain support
func (r *RestErr) Unwrap() error {
	return r.Wrapped
}

// WithCause wraps an underlying error
func (r *RestErr) WithCause(err error) *RestErr {
	r.Wrapped = err
	return r
}

// IsClientError returns true if the error is a 4xx client error
func (r *RestErr) IsClientError() bool {
	return r.Code >= 400 && r.Code < 500
}

// IsServerError returns true if the error is a 5xx server error
func (r *RestErr) IsServerError() bool {
	return r.Code >= 500 && r.Code < 600
}

// IsNotFound returns true if the error is a 404 Not Found
func (r *RestErr) IsNotFound() bool {
	return r.Code == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func (r *RestErr) IsUnauthorized() bool {
	return r.Code == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden
func (r *RestErr) IsForbidden() bool {
	return r.Code == http.StatusForbidden
}

func NewRestErr(message, err string, code int, causes []Causes) *RestErr {
	return &RestErr{
		Message:   message,
		Err:       err,
		Code:      code,
		Causes:    causes,
		Timestamp: time.Now(),
	}
}

// NewRestErrFromError converts a standard Go error to a RestErr
// Defaults to 500 Internal Server Error
func NewRestErrFromError(err error) *RestErr {
	if err == nil {
		return nil
	}

	// Check if it's already a RestErr
	var restErr *RestErr
	if errors.As(err, &restErr) {
		return restErr
	}

	// Default to internal server error
	return &RestErr{
		Message:   "An unexpected error occurred",
		Err:       "internal server error",
		Code:      http.StatusInternalServerError,
		Wrapped:   err,
		Timestamp: time.Now(),
	}
}

// ParseError attempts to extract a RestErr from an error chain
func ParseError(err error) (*RestErr, bool) {
	if err == nil {
		return nil, false
	}
	var restErr *RestErr
	if errors.As(err, &restErr) {
		return restErr, true
	}
	return nil, false
}

func NewBadRequestError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "bad request",
		Code:      http.StatusBadRequest,
		Timestamp: time.Now(),
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message:   message,
		Err:       "bad request",
		Code:      http.StatusBadRequest,
		Causes:    causes,
		Timestamp: time.Now(),
	}
}

func NewInternalServerError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "internal server error",
		Code:      http.StatusInternalServerError,
		Timestamp: time.Now(),
	}
}

func NewNotFoundError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "not found",
		Code:      http.StatusNotFound,
		Timestamp: time.Now(),
	}
}

func NewForbiddenError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "forbidden",
		Code:      http.StatusForbidden,
		Timestamp: time.Now(),
	}
}

func NewUnauthorizedError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "unauthorized",
		Code:      http.StatusUnauthorized,
		Timestamp: time.Now(),
	}
}

func NewBadGatewayError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "bad gateway",
		Code:      http.StatusBadGateway,
		Timestamp: time.Now(),
	}
}

func NewConflictError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "conflict",
		Code:      http.StatusConflict,
		Timestamp: time.Now(),
	}
}

func NewUnprocessableEntityError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message:   message,
		Err:       "unprocessable entity",
		Code:      http.StatusUnprocessableEntity,
		Causes:    causes,
		Timestamp: time.Now(),
	}
}

func NewTooManyRequestsError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "too many requests",
		Code:      http.StatusTooManyRequests,
		Timestamp: time.Now(),
	}
}

func NewServiceUnavailableError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "service unavailable",
		Code:      http.StatusServiceUnavailable,
		Timestamp: time.Now(),
	}
}

func NewGatewayTimeoutError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "gateway timeout",
		Code:      http.StatusGatewayTimeout,
		Timestamp: time.Now(),
	}
}

func NewPreconditionFailedError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "precondition failed",
		Code:      http.StatusPreconditionFailed,
		Timestamp: time.Now(),
	}
}

func NewNotAcceptableError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "not acceptable",
		Code:      http.StatusNotAcceptable,
		Timestamp: time.Now(),
	}
}

func NewLengthRequiredError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "length required",
		Code:      http.StatusLengthRequired,
		Timestamp: time.Now(),
	}
}

func NewUnsupportedMediaTypeError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "unsupported media type",
		Code:      http.StatusUnsupportedMediaType,
		Timestamp: time.Now(),
	}
}

func NewExpectationFailedError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "expectation failed",
		Code:      http.StatusExpectationFailed,
		Timestamp: time.Now(),
	}
}

func NewConflictValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message:   message,
		Err:       "conflict",
		Code:      http.StatusConflict,
		Causes:    causes,
		Timestamp: time.Now(),
	}
}

func NewRequestTimeoutError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "request timeout",
		Code:      http.StatusRequestTimeout,
		Timestamp: time.Now(),
	}
}

func NewHttpVersionNotSupportedError(message string, args ...any) *RestErr {
	return &RestErr{
		Message:   fmt.Sprintf(message, args...),
		Err:       "http version not supported",
		Code:      http.StatusHTTPVersionNotSupported,
		Timestamp: time.Now(),
	}
}

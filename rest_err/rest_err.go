package rest_err

import "net/http"

type RestErr struct {
	Message string   `json:"message" example:"invalid request parameters"` // Human readable message
	Err     string   `json:"error" example:"bad request"`
	Code    int      `json:"code" example:"400"` // HTTP status code
	Causes  []Causes `json:"causes"`             // Detailed error causes, most common for json field validation errors
}

type Causes struct {
	Field   string `json:"field" example:"email"`                   // Field or parameter that caused the error
	Message string `json:"message" example:"invalid email address"` // Description of the cause
}

func (r *RestErr) Error() string {
	return r.Message
}

func NewRestErr(message, err string, code int, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     err,
		Code:    code,
		Causes:  causes,
	}
}

func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad request",
		Code:    http.StatusBadRequest,
	}
}

func NewBadRequestValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad request",
		Code:    http.StatusBadRequest,
		Causes:  causes,
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "internal server error",
		Code:    http.StatusInternalServerError,
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not found",
		Code:    http.StatusNotFound,
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "forbidden",
		Code:    http.StatusForbidden,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unauthorized",
		Code:    http.StatusUnauthorized,
	}
}

func NewBadGatewayError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "bad gateway",
		Code:    http.StatusBadGateway,
	}
}

func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
	}
}

func NewUnprocessableEntityError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unprocessable entity",
		Code:    http.StatusUnprocessableEntity,
		Causes:  causes,
	}
}

func NewTooManyRequestsError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "too many requests",
		Code:    http.StatusTooManyRequests,
	}
}

func NewServiceUnavailableError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "service unavailable",
		Code:    http.StatusServiceUnavailable,
	}
}

func NewGatewayTimeoutError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "gateway timeout",
		Code:    http.StatusGatewayTimeout,
	}
}

func NewPreconditionFailedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "precondition failed",
		Code:    http.StatusPreconditionFailed,
	}
}

func NewNotAcceptableError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "not acceptable",
		Code:    http.StatusNotAcceptable,
	}
}

func NewLengthRequiredError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "length required",
		Code:    http.StatusLengthRequired,
	}
}

func NewUnsupportedMediaTypeError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "unsupported media type",
		Code:    http.StatusUnsupportedMediaType,
	}
}

func NewExpectationFailedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "expectation failed",
		Code:    http.StatusExpectationFailed,
	}
}

func NewConflictValidationError(message string, causes []Causes) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "conflict",
		Code:    http.StatusConflict,
		Causes:  causes,
	}
}

func NewRequestTimeoutError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "request timeout",
		Code:    http.StatusRequestTimeout,
	}
}

func NewHttpVersionNotSupportedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Err:     "http version not supported",
		Code:    http.StatusHTTPVersionNotSupported,
	}
}

package rest_err

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewRestErr(t *testing.T) {
	causes := []Causes{{Field: "email", Message: "invalid format"}}
	err := NewRestErr("test message", "test error", 400, causes)

	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}
	if err.Err != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", err.Err)
	}
	if err.Code != 400 {
		t.Errorf("Expected code 400, got %d", err.Code)
	}
	if len(err.Causes) != 1 {
		t.Errorf("Expected 1 cause, got %d", len(err.Causes))
	}
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestRestErr_Error(t *testing.T) {
	t.Run("without wrapped error", func(t *testing.T) {
		err := NewBadRequestError("test message")
		if err.Error() != "test message" {
			t.Errorf("Expected 'test message', got '%s'", err.Error())
		}
	})

	t.Run("with wrapped error", func(t *testing.T) {
		wrappedErr := errors.New("underlying error")
		err := NewInternalServerError("server error").WithCause(wrappedErr)
		expected := "server error: underlying error"
		if err.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestRestErr_Unwrap(t *testing.T) {
	wrappedErr := errors.New("underlying error")
	err := NewInternalServerError("server error").WithCause(wrappedErr)

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Error("Expected unwrapped error to match original wrapped error")
	}
}

func TestRestErr_WithCause(t *testing.T) {
	originalErr := errors.New("database connection failed")
	restErr := NewInternalServerError("failed to fetch data").WithCause(originalErr)

	if restErr.Wrapped != originalErr {
		t.Error("Expected wrapped error to be set")
	}

	// Test error chain
	if !errors.Is(restErr, originalErr) {
		t.Error("Expected error chain to work with errors.Is")
	}
}

func TestRestErr_IsClientError(t *testing.T) {
	tests := []struct {
		name     string
		err      *RestErr
		expected bool
	}{
		{"400 Bad Request", NewBadRequestError("test"), true},
		{"404 Not Found", NewNotFoundError("test"), true},
		{"500 Internal Server Error", NewInternalServerError("test"), false},
		{"503 Service Unavailable", NewServiceUnavailableError("test"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.IsClientError() != tt.expected {
				t.Errorf("Expected IsClientError() to be %v for %s", tt.expected, tt.name)
			}
		})
	}
}

func TestRestErr_IsServerError(t *testing.T) {
	tests := []struct {
		name     string
		err      *RestErr
		expected bool
	}{
		{"400 Bad Request", NewBadRequestError("test"), false},
		{"404 Not Found", NewNotFoundError("test"), false},
		{"500 Internal Server Error", NewInternalServerError("test"), true},
		{"503 Service Unavailable", NewServiceUnavailableError("test"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.IsServerError() != tt.expected {
				t.Errorf("Expected IsServerError() to be %v for %s", tt.expected, tt.name)
			}
		})
	}
}

func TestRestErr_IsNotFound(t *testing.T) {
	notFoundErr := NewNotFoundError("resource not found")
	badRequestErr := NewBadRequestError("bad request")

	if !notFoundErr.IsNotFound() {
		t.Error("Expected IsNotFound() to return true for 404 error")
	}
	if badRequestErr.IsNotFound() {
		t.Error("Expected IsNotFound() to return false for non-404 error")
	}
}

func TestRestErr_IsUnauthorized(t *testing.T) {
	unauthorizedErr := NewUnauthorizedError("not authorized")
	badRequestErr := NewBadRequestError("bad request")

	if !unauthorizedErr.IsUnauthorized() {
		t.Error("Expected IsUnauthorized() to return true for 401 error")
	}
	if badRequestErr.IsUnauthorized() {
		t.Error("Expected IsUnauthorized() to return false for non-401 error")
	}
}

func TestRestErr_IsForbidden(t *testing.T) {
	forbiddenErr := NewForbiddenError("forbidden")
	badRequestErr := NewBadRequestError("bad request")

	if !forbiddenErr.IsForbidden() {
		t.Error("Expected IsForbidden() to return true for 403 error")
	}
	if badRequestErr.IsForbidden() {
		t.Error("Expected IsForbidden() to return false for non-403 error")
	}
}

func TestNewRestErrFromError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := NewRestErrFromError(nil)
		if err != nil {
			t.Error("Expected nil for nil input")
		}
	})

	t.Run("standard error", func(t *testing.T) {
		standardErr := errors.New("something went wrong")
		restErr := NewRestErrFromError(standardErr)

		if restErr.Code != http.StatusInternalServerError {
			t.Errorf("Expected code 500, got %d", restErr.Code)
		}
		if restErr.Wrapped != standardErr {
			t.Error("Expected wrapped error to be set")
		}
		if !errors.Is(restErr, standardErr) {
			t.Error("Expected error chain to work")
		}
	})

	t.Run("already a RestErr", func(t *testing.T) {
		originalRestErr := NewNotFoundError("user not found")
		wrappedErr := fmt.Errorf("wrapped: %w", originalRestErr)
		result := NewRestErrFromError(wrappedErr)

		if result.Code != http.StatusNotFound {
			t.Errorf("Expected code 404, got %d", result.Code)
		}
		if result.Message != "user not found" {
			t.Errorf("Expected original message, got '%s'", result.Message)
		}
	})
}

func TestParseError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		restErr, ok := ParseError(nil)
		if ok || restErr != nil {
			t.Error("Expected false and nil for nil input")
		}
	})

	t.Run("standard error", func(t *testing.T) {
		standardErr := errors.New("standard error")
		restErr, ok := ParseError(standardErr)
		if ok || restErr != nil {
			t.Error("Expected false and nil for standard error")
		}
	})

	t.Run("direct RestErr", func(t *testing.T) {
		originalErr := NewBadRequestError("bad request")
		restErr, ok := ParseError(originalErr)
		if !ok || restErr != originalErr {
			t.Error("Expected true and the original RestErr")
		}
	})

	t.Run("wrapped RestErr", func(t *testing.T) {
		originalErr := NewNotFoundError("not found")
		wrappedErr := fmt.Errorf("context: %w", originalErr)
		restErr, ok := ParseError(wrappedErr)
		if !ok {
			t.Error("Expected true for wrapped RestErr")
		}
		if restErr.Code != http.StatusNotFound {
			t.Errorf("Expected code 404, got %d", restErr.Code)
		}
	})
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name         string
		constructor  func() *RestErr
		expectedCode int
		expectedErr  string
	}{
		{"BadRequest", func() *RestErr { return NewBadRequestError("bad request") }, http.StatusBadRequest, "bad request"},
		{"InternalServer", func() *RestErr { return NewInternalServerError("server error") }, http.StatusInternalServerError, "internal server error"},
		{"NotFound", func() *RestErr { return NewNotFoundError("not found") }, http.StatusNotFound, "not found"},
		{"Forbidden", func() *RestErr { return NewForbiddenError("forbidden") }, http.StatusForbidden, "forbidden"},
		{"Unauthorized", func() *RestErr { return NewUnauthorizedError("unauthorized") }, http.StatusUnauthorized, "unauthorized"},
		{"BadGateway", func() *RestErr { return NewBadGatewayError("bad gateway") }, http.StatusBadGateway, "bad gateway"},
		{"Conflict", func() *RestErr { return NewConflictError("conflict") }, http.StatusConflict, "conflict"},
		{"TooManyRequests", func() *RestErr { return NewTooManyRequestsError("too many") }, http.StatusTooManyRequests, "too many requests"},
		{"ServiceUnavailable", func() *RestErr { return NewServiceUnavailableError("unavailable") }, http.StatusServiceUnavailable, "service unavailable"},
		{"GatewayTimeout", func() *RestErr { return NewGatewayTimeoutError("timeout") }, http.StatusGatewayTimeout, "gateway timeout"},
		{"PreconditionFailed", func() *RestErr { return NewPreconditionFailedError("failed") }, http.StatusPreconditionFailed, "precondition failed"},
		{"NotAcceptable", func() *RestErr { return NewNotAcceptableError("not acceptable") }, http.StatusNotAcceptable, "not acceptable"},
		{"LengthRequired", func() *RestErr { return NewLengthRequiredError("length required") }, http.StatusLengthRequired, "length required"},
		{"UnsupportedMediaType", func() *RestErr { return NewUnsupportedMediaTypeError("unsupported") }, http.StatusUnsupportedMediaType, "unsupported media type"},
		{"ExpectationFailed", func() *RestErr { return NewExpectationFailedError("expectation") }, http.StatusExpectationFailed, "expectation failed"},
		{"RequestTimeout", func() *RestErr { return NewRequestTimeoutError("timeout") }, http.StatusRequestTimeout, "request timeout"},
		{"HTTPVersionNotSupported", func() *RestErr { return NewHttpVersionNotSupportedError("not supported") }, http.StatusHTTPVersionNotSupported, "http version not supported"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.constructor()
			if err.Code != tt.expectedCode {
				t.Errorf("Expected code %d, got %d", tt.expectedCode, err.Code)
			}
			if err.Err != tt.expectedErr {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Err)
			}
			if err.Timestamp.IsZero() {
				t.Error("Expected timestamp to be set")
			}
			if time.Since(err.Timestamp) > time.Second {
				t.Error("Expected timestamp to be recent")
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	causes := []Causes{
		{Field: "email", Message: "invalid email format"},
		{Field: "password", Message: "too short"},
	}

	t.Run("BadRequestValidation", func(t *testing.T) {
		err := NewBadRequestValidationError("validation failed", causes)
		if err.Code != http.StatusBadRequest {
			t.Errorf("Expected code 400, got %d", err.Code)
		}
		if len(err.Causes) != 2 {
			t.Errorf("Expected 2 causes, got %d", len(err.Causes))
		}
	})

	t.Run("UnprocessableEntity", func(t *testing.T) {
		err := NewUnprocessableEntityError("cannot process", causes)
		if err.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected code 422, got %d", err.Code)
		}
		if len(err.Causes) != 2 {
			t.Errorf("Expected 2 causes, got %d", len(err.Causes))
		}
	})

	t.Run("ConflictValidation", func(t *testing.T) {
		err := NewConflictValidationError("conflict", causes)
		if err.Code != http.StatusConflict {
			t.Errorf("Expected code 409, got %d", err.Code)
		}
		if len(err.Causes) != 2 {
			t.Errorf("Expected 2 causes, got %d", len(err.Causes))
		}
	})
}

func TestErrorFormatting(t *testing.T) {
	t.Run("with formatting args", func(t *testing.T) {
		err := NewBadRequestError("user %s not found with id %d", "john", 123)
		expected := "user john not found with id 123"
		if err.Message != expected {
			t.Errorf("Expected message '%s', got '%s'", expected, err.Message)
		}
	})

	t.Run("without formatting args", func(t *testing.T) {
		err := NewBadRequestError("simple message")
		if err.Message != "simple message" {
			t.Errorf("Expected message 'simple message', got '%s'", err.Message)
		}
	})
}

func TestErrorChaining(t *testing.T) {
	t.Run("errors.Is with wrapped error", func(t *testing.T) {
		originalErr := errors.New("database error")
		restErr := NewInternalServerError("failed to save").WithCause(originalErr)

		if !errors.Is(restErr, originalErr) {
			t.Error("Expected errors.Is to work with wrapped error")
		}
	})

	t.Run("errors.As with RestErr", func(t *testing.T) {
		restErr := NewNotFoundError("not found")
		wrappedErr := fmt.Errorf("wrapped: %w", restErr)

		var targetErr *RestErr
		if !errors.As(wrappedErr, &targetErr) {
			t.Error("Expected errors.As to extract RestErr")
		}
		if targetErr.Code != http.StatusNotFound {
			t.Errorf("Expected code 404, got %d", targetErr.Code)
		}
	})

	t.Run("multiple levels of wrapping", func(t *testing.T) {
		baseErr := errors.New("base error")
		restErr := NewInternalServerError("server error").WithCause(baseErr)
		wrappedErr := fmt.Errorf("layer 1: %w", restErr)
		finalErr := fmt.Errorf("layer 2: %w", wrappedErr)

		var targetErr *RestErr
		if !errors.As(finalErr, &targetErr) {
			t.Error("Expected errors.As to extract RestErr through multiple wraps")
		}
		if !errors.Is(finalErr, baseErr) {
			t.Error("Expected errors.Is to find base error through multiple wraps")
		}
	})
}

func BenchmarkNewBadRequestError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewBadRequestError("test error")
	}
}

func BenchmarkNewRestErrFromError(b *testing.B) {
	err := errors.New("test error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewRestErrFromError(err)
	}
}

func BenchmarkWithCause(b *testing.B) {
	baseErr := errors.New("base error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewInternalServerError("server error").WithCause(baseErr)
	}
}

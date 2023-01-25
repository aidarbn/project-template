package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"net/http"
)

// HTTPError is a general error returned by REST API
type HTTPError struct {
	// Example: 500
	Code int `json:"code"`

	// Example: Unexpected internal server error
	Description string `json:"description"`

	// Wrapped error
	Err error `json:"-"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Description)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// WithError wraps the internal error
func (e *HTTPError) WithError(err error) *HTTPError {
	e.Err = err
	return e
}

// NewHTTPError returns REST error with code and message
func NewHTTPError(code int, format string, args ...any) *HTTPError {
	return &HTTPError{
		Code:        code,
		Description: fmt.Sprintf(format, args...),
	}
}

// BadRequestErrorf returns REST error with 400 status code and message.
func BadRequestErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, format, args...)
}

// InternalServerErrorf returns REST error with 500 status code and message.
func InternalServerErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, format, args...)
}

// NotFoundErrorf returns REST error with 404 status code and message.
func NotFoundErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusNotFound, format, args...)
}

// UnauthorizedErrorf returns REST error with 401 status code and message.
func UnauthorizedErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, format, args...)
}

// ForbiddenErrorf returns REST error with 403 status code and message.
func ForbiddenErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusForbidden, format, args...)
}

// ConflictErrorf returns REST error with 409 status code and message.
func ConflictErrorf(format string, args ...any) *HTTPError {
	return NewHTTPError(http.StatusConflict, format, args...)
}

// APIHandler is type that extends standard http handler func with error.
type APIHandler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP calls API handler and handles all errors and panics.
func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		WriteError(w, r, err)
	}
}

// APIHandlerFunc gets APIHandler and returns standard http.HandlerFunc
func APIHandlerFunc(fn APIHandler) http.HandlerFunc {
	return fn.ServeHTTP
}

// MiddlewareHandler is a type that extends standard middleware func with response, request and error.
type MiddlewareHandler func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error)

// Middleware returns new handler processed by middleware.
func (m MiddlewareHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if w, r, err = m(w, r); err != nil {
			WriteError(w, r, err)

			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareHandlerFunc gets MiddlewareHandler and returns standard middleware
// function.
func MiddlewareHandlerFunc(fn MiddlewareHandler) func(http.Handler) http.Handler {
	return fn.Middleware
}

// WriteError logs detailed message and sends encoded error to the client.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	log := getLogEntry(r)
	var apiErr *HTTPError
	if errors.As(err, &apiErr) {
		entry := log
		if apiErr.Err != nil {
			entry = entry.WithError(apiErr.Err)
		}
		if apiErr.Err != nil {
			entry.WithError(apiErr.Err).Error(apiErr)
		} else {
			entry.Error(apiErr)
		}
		if sendErr := WriteJSON(w, apiErr, apiErr.Code); sendErr != nil {
			log.WithError(sendErr).Error(err)
		}
	} else {
		sentryHub := sentry.GetHubFromContext(r.Context())
		if sentryHub != nil {
			sentryHub.CaptureException(err)
		}
		apiErr = InternalServerErrorf("Internal Server Error")
		log.WithError(err).Error(apiErr)
		if sendErr := WriteJSON(w, apiErr, apiErr.Code); sendErr != nil {
			log.WithError(sendErr).Error(err)
		}
	}
}

// WriteJSON sends value v to response writer w as JSON.
func WriteJSON(w http.ResponseWriter, v any, statusCode int) error {
	js, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("can't encode %v in JSON: %w", v, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err = w.Write(js); err != nil {
		return fmt.Errorf("can't write response: %w", err)
	}
	return nil
}

func WriteNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

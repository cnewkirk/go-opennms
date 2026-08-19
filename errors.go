package opennms

import (
	"errors"
	"fmt"
)

// Sentinel errors for HTTP error classes. Match them with errors.Is
// against errors returned by any Client method:
//
//	_, err := client.GetNode(ctx, 99999)
//	if errors.Is(err, opennms.ErrNotFound) { ... }
var (
	// ErrBadRequest matches 400 Bad Request responses.
	ErrBadRequest = errors.New("opennms: bad request")
	// ErrAuthentication matches 401 Unauthorized responses.
	ErrAuthentication = errors.New("opennms: authentication failed")
	// ErrForbidden matches 403 Forbidden responses.
	ErrForbidden = errors.New("opennms: forbidden")
	// ErrNotFound matches 404 Not Found responses.
	ErrNotFound = errors.New("opennms: not found")
	// ErrConflict matches 409 Conflict responses.
	ErrConflict = errors.New("opennms: conflict")
	// ErrServer matches 5xx server error responses.
	ErrServer = errors.New("opennms: server error")
)

// APIError is returned when the OpenNMS server answers with a non-2xx
// status code. Use errors.As to inspect the status code and body, and
// errors.Is with the package sentinel errors to match error classes.
type APIError struct {
	// StatusCode is the HTTP status code returned by the server.
	StatusCode int
	// Status is the full HTTP status line, e.g. "404 Not Found".
	Status string
	// URL is the request URL that produced the error.
	URL string
	// Body is the raw response body, if any.
	Body string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("opennms: HTTP %s for %s", e.Status, e.URL)
}

// Is reports whether the error matches one of the package sentinel
// errors based on its status code.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrBadRequest:
		return e.StatusCode == 400
	case ErrAuthentication:
		return e.StatusCode == 401
	case ErrForbidden:
		return e.StatusCode == 403
	case ErrNotFound:
		return e.StatusCode == 404
	case ErrConflict:
		return e.StatusCode == 409
	case ErrServer:
		return e.StatusCode >= 500 && e.StatusCode < 600
	}
	return false
}

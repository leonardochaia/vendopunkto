package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

func InvalidRequestParams(err error) *APIError {
	return &APIError{
		Title:  "Supplied parameters were invalid",
		Status: 400,
		Detail: err.Error(),
	}
}

func ResourceNotFound() *APIError {
	return &APIError{
		Title:  "Requested resource was not found",
		Status: 404,
	}
}

func InternalServerError(err error) *APIError {
	return &APIError{
		Title:  "Internal Server Error",
		Status: 500,
		Detail: err.Error(),
	}
}

// Based on https://tools.ietf.org/html/rfc7807#section-3.1
type APIError struct {
	Type   string `json:"type,omitempty"`   // A URI reference that identifies the problem type.
	Title  string `json:"title"`            // A short, human-readable summary of the problem type.
	Status int    `json:"status"`           // This is for client's convinience.
	Detail string `json:"detail,omitempty"` //  A human-readable explanation specific to this occurrence of the problem.
}

func (e *APIError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)
	return nil
}

type HandlerErrorAwareFunc func(http.ResponseWriter, *http.Request) *APIError

func WrapHandler(handler HandlerErrorAwareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			render.Render(w, r, err)
		}
	}
}

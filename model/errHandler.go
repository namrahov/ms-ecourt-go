package model

import (
	"encoding/json"
	"net/http"
)

type CommonError struct {
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type FileNotFoundError CommonError

func (e *FileNotFoundError) Error() string {
	return e.Message
}

func NewFileNotFoundError(message string) error {
	return &FileNotFoundError{
		Message: message,
		Code:    "file.not.found",
	}
}

func ErrorHandler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the final handler, and deal with errors
		err := h(w, r)
		if err != nil {
			switch e := err.(type) {
			case *FileNotFoundError:
				// We can retrieve the status here and write out a specific
				// HTTP status code.
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(e)
			default:
				// Any error types we don't specifically look out for default
				// to serving a HTTP 500
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}
	}
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

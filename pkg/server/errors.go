package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrRequiredContentTypeJSON = errors.New("header Content-Type must be application/json")
	ErrInternal                = errors.New(http.StatusText(http.StatusInternalServerError))
	ErrEmptyRequestBody        = errors.New("request body must not be empty")
)

type ErrorResponse struct {
	Err string `json:"error"`
}

func writeErrorAsJSON(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Err: err.Error(),
	})
}

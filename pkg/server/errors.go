package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrRequiredContentTypeJSON = errors.New("header Content-Type is not application/json")
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

package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrRequiredContentTypeJSON = errors.New("Content-Type header is not JSON")
)

func writeErrorAsJSON(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

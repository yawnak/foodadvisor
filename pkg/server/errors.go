package server

import (
	"encoding/json"
	"net/http"
)

func writeErrorAsJSON(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

package server

import (
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil/header"
)



// validates that request content type is set to JSON
func validateContentJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("validating content-type")
		fmt.Println("content type: ", r.Header.Get("Content-Type"))
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			writeErrorAsJSON(w, http.StatusUnsupportedMediaType, ErrRequiredContentTypeJSON)
			return
		}
		next.ServeHTTP(w, r)
	}
}

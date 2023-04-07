package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
)

// validates that request content type is set to JSON
func ValidateContentJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			exception.WriteErrorAsJSON(w, http.StatusUnsupportedMediaType, exception.ErrRequiredContentTypeJSON)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ConfirmPermissions(permissions ...domain.Permission) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := domain.RoleFromContext(r.Context())
			if !ok {
				exception.WriteErrorAsJSON(w, http.StatusUnauthorized, fmt.Errorf("role is missing"))
				return
			}

			for _, p := range permissions {
				_, ok := role.Permissions[p]
				if !ok {
					exception.WriteErrorAsJSON(w, http.StatusUnauthorized, fmt.Errorf("permissions %s is missing", p))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ToMiddleware(hf http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hf(w, r)
		})
	}
}

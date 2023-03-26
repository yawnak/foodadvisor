package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/yawnak/foodadvisor/pkg/domain"
)

// validates that request content type is set to JSON
func validateContentJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			writeErrorAsJSON(w, http.StatusUnsupportedMediaType, ErrRequiredContentTypeJSON)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (srv *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.authTokenToContext(next.ServeHTTP)(w, r)
	})
}

func retrieveUserId(ctx context.Context) (int32, bool) {
	id, ok := ctx.Value(keyUserId).(int32)
	return id, ok
}

func retrieveRole(ctx context.Context) (*domain.Role, bool) {
	role, ok := ctx.Value(keyRole).(*domain.Role)
	return role, ok
}

func confirmPermissions(permissions ...domain.Permission) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := retrieveRole(r.Context())
			if !ok {
				writeErrorAsJSON(w, http.StatusUnauthorized, fmt.Errorf("role is missing"))
				return
			}

			for _, p := range permissions {
				_, ok := role.Permissions[p]
				if !ok {
					writeErrorAsJSON(w, http.StatusUnauthorized, fmt.Errorf("permissions %s is missing", p))
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/xorcare/pointer"
	"github.com/yawnak/foodadvisor/pkg/bind"
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

// return http.StatusOK if err = nil
// return http.StatusInternalServerError if error is unknown
// other - if error is known (mostly http.StatusBadRequest)
func bindingErrorToHTTPStatus(err error) int {
	if err != nil {
		switch {
		case errors.Is(err, bind.ErrBadFormat):
			return http.StatusBadRequest
		case errors.Is(err, bind.ErrBodyTooLarge):
			return http.StatusBadRequest
		case errors.Is(err, bind.ErrEmptyBody):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrSyntax{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnmarshalType{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnknownField{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnknown{})):
			return http.StatusInternalServerError
		}
	}
	return http.StatusOK
}

func ToMiddleware(hf http.HandlerFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hf(w, r)
		})
	}
}

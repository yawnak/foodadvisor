package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/golang/gddo/httputil/header"
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

			case errors.Is(err, domain.ErrInvalidSigningMethod):
				writeErrorAsJSON(w, http.StatusUnauthorized, err)
				return
			default:
				log.Println(err)
				writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("error validating auth token cookie"))
				return
			}
		}
		next(w, r.WithContext(context.WithValue(r.Context(), keyUserId, id)))
	}
}

func retrieveUserId(ctx context.Context) (int32, bool) {
	id, ok := ctx.Value(keyUserId).(int32)
	return id, ok
}

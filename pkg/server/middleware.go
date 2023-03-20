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

func (srv *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				writeErrorAsJSON(w, http.StatusUnauthorized, errors.New("cookie with auth token is not present"))
				return
			}
		}
		id, err := srv.app.ParseToken(r.Context(), cookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrBadToken):
				writeErrorAsJSON(w, http.StatusUnauthorized, err)
				return
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

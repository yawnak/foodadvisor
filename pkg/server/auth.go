package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yawnak/foodadvisor/pkg/domain"
)

const (
	authCookieName = "foodAdvisorAuthToken"
)

func (srv *Server) signup(w http.ResponseWriter, r *http.Request) {
	var err error
	//unmarshaling request to domain.User
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) //set maximum size of request body to 1mb
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //will not allow sending unknown fields in request
	var user domain.User
	err = dec.Decode(&user)
	//processing errors
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			err = fmt.Errorf("request body contains bad JSON (at position %d)", syntaxError.Offset)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.Is(err, io.ErrUnexpectedEOF):
			err = errors.New("request body contains badly-formed JSON")
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.As(err, &unmarshalTypeError):
			err = fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			err = fmt.Errorf("request body contains unknown field %s", fieldName)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.Is(err, io.EOF):
			writeErrorAsJSON(w, http.StatusBadRequest, ErrEmptyRequestBody)
			return
		default:
			log.Print(err.Error())
			writeErrorAsJSON(w, http.StatusInternalServerError, ErrInternal)
			return
		}
	}
	id, err := srv.app.CreateUser(r.Context(), &user)
	if err != nil {
		fmt.Fprintf(w, "error creating user: %s", err)
		return
	}
	token, err := srv.app.GenerateToken(r.Context(), user.Username, user.Password)
	if err != nil {
		log.Println("error when generating token after signup:", err)
		writeErrorAsJSON(w, http.StatusUnauthorized, errors.New("error creating auth token"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   token,
		Expires: time.Now().Add(domain.TokenTTL),
	})

	fmt.Fprint(w, id)
}

func (srv *Server) login(w http.ResponseWriter, r *http.Request) {
	var err error
	//unmarshaling request to domain.User
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) //set maximum size of request body to 1mb
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //will not allow sending unknown fields in request
	var credentials requestLogin
	err = dec.Decode(&credentials)
	//processing errors
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			err = fmt.Errorf("request body contains bad JSON (at position %d)", syntaxError.Offset)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.Is(err, io.ErrUnexpectedEOF):
			err = errors.New("request body contains badly-formed JSON")
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.As(err, &unmarshalTypeError):
			err = fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			err = fmt.Errorf("request body contains unknown field %s", fieldName)
			writeErrorAsJSON(w, http.StatusBadRequest, err)
			return
		case errors.Is(err, io.EOF):
			writeErrorAsJSON(w, http.StatusBadRequest, ErrEmptyRequestBody)
			return
		default:
			log.Print(err.Error())
			writeErrorAsJSON(w, http.StatusInternalServerError, ErrInternal)
			return
		}
	}
	token, err := srv.app.GenerateToken(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWrongCredentials):
			writeErrorAsJSON(w, http.StatusUnauthorized, err)
			return
		default:
			log.Println("error generating token while login:", err)
			writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("error generating auth token"))
			return
		}
	}
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   token,
		Expires: time.Now().Add(domain.TokenTTL),
	})
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) authTokenToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				writeErrorAsJSON(w, http.StatusUnauthorized, errors.New("cookie with auth token is not present"))
				return
			}
		}
		id, role, err := srv.app.ParseTokenWithRole(r.Context(), cookie.Value)
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
		ctxnext := context.WithValue(r.Context(), keyUserId, id)
		ctxnext = context.WithValue(ctxnext, keyRole, role)
		next(w, r.WithContext(ctxnext))
	}
}

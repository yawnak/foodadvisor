package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/bind"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
)

const (
	authCookieName = "foodAdvisorAuthToken"
)

type userIdCtxKey struct{}

func userIdToContext(ctx context.Context, userid int32) context.Context {
	return context.WithValue(ctx, userIdCtxKey{}, userid)
}

func userIdFromContext(ctx context.Context) (int32, bool) {
	uid, ok := ctx.Value(userIdCtxKey{}).(int32)
	if !ok {
		return 0, ok
	}
	return uid, ok
}

func (srv *Server) signup(w http.ResponseWriter, r *http.Request) {
	var err error
	var user domain.User
	//unmarshaling request to domain.User
	binder := bind.JSONBinder{}
	var userReq requestSignup
	err = binder.Bind(&userReq, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})
	//conver error to http StatusCode
	status := bindingErrorToHTTPStatus(err)
	switch status {
	case http.StatusOK: //if ok continue
	case http.StatusInternalServerError: //if internal server error log it
		log.Printf("unknown error while binding requestSignup: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unknown error while parsing request body"))
		return
	default: //other errors are printed to user
		exception.WriteErrorAsJSON(w, status, err)
		return
	}
	//struct validation
	err = srv.validate.Struct(userReq)
	if err != nil {
		exception.WriteErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}
	user = domain.User(userReq) //converting requestSignup to domain.User
	//creating user
	id, err := srv.app.CreateUser(r.Context(), &user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDuplicateResourse): //if user already exists
			exception.WriteErrorAsJSON(w, http.StatusBadRequest, domain.ErrDuplicateResourse)
			return
		case errors.Is(err, domain.ErrPasswordTooLong): //if password is too long (this shouldn't happen because of validation)
			exception.WriteErrorAsJSON(w, http.StatusBadRequest, domain.ErrPasswordTooLong)
			return
		}
		//everything else is logged and user recieves internal server error
		log.Printf("unexpected error creating user: %s", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("unexpected error while creating user"))
		return
	}
	//token generation
	token, err := srv.app.GenerateToken(r.Context(), user.Username, user.Password)
	if err != nil {
		log.Println("error when generating token after signup:", err)
		exception.WriteErrorAsJSON(w, http.StatusUnauthorized, errors.New("error creating auth token"))
		return
	}
	//setting cookie
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   token,
		Expires: time.Now().Add(domain.TokenTTL),
	})

	writeSuccess(w, responseSignup{
		UserId: id,
		responseSuccess: responseSuccess{
			HTTPStatusCode: http.StatusOK,
			SuccessMessage: "ok",
		},
	})
}

func (srv *Server) login(w http.ResponseWriter, r *http.Request) {
	var err error
	//unmarshaling request to domain.User
	var credentials requestLogin
	binder := bind.JSONBinder{}
	err = binder.Bind(&credentials, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})
	//processing errors
	switch status := bindingErrorToHTTPStatus(err); status {
	case http.StatusOK:
	case http.StatusInternalServerError:
		log.Printf("unexpected error binding requestLogin: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unexpected error while parsing request body"))
		return
	default:
		exception.WriteErrorAsJSON(w, status, err)
		return
	}

	//struct validation
	err = srv.validate.Struct(credentials)
	if err != nil {
		exception.WriteErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	//token generation
	token, err := srv.app.GenerateToken(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWrongCredentials):
			exception.WriteErrorAsJSON(w, http.StatusUnauthorized, err)
			return
		default:
			log.Println("error generating token while login:", err)
			exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("error generating auth token"))
			return
		}
	}
	//setting cookie with auth token
	http.SetCookie(w, &http.Cookie{
		Name:    authCookieName,
		Value:   token,
		Expires: time.Now().Add(domain.TokenTTL),
	})

	writeSuccessOK(w, "ok")
}

func (srv *Server) authTokenToContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookieName)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				exception.WriteErrorAsJSON(w, http.StatusUnauthorized, errors.New("cookie with auth token is not present"))
				return
			}
		}
		id, role, err := srv.app.ParseTokenWithRole(r.Context(), cookie.Value)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrBadToken):
				exception.WriteErrorAsJSON(w, http.StatusUnauthorized, err)
				return
			case errors.Is(err, domain.ErrInvalidSigningMethod):
				exception.WriteErrorAsJSON(w, http.StatusUnauthorized, err)
				return
			default:
				log.Println(err)
				exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("error validating auth token cookie"))
				return
			}
		}
		r = r.WithContext(context.WithValue(r.Context(), keyUserId, id))
		r = r.WithContext(role.ToContext(r.Context()))
		next(w, r)
	}
}

func (srv *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.authTokenToContext(next.ServeHTTP)(w, r)
	})
}

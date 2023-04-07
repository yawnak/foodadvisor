package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/bind"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
)

type key int

const (
	keyUserId key = iota
	keyRole
)

type Successer interface {
	Success() string
	Status() int
}

func writeSuccess(w http.ResponseWriter, success Successer) {
	w.WriteHeader(success.Status())
	_ = json.NewEncoder(w).Encode(success)
}

func writeSuccessOK(w http.ResponseWriter, message string) {
	writeSuccess(w, responseSuccess{SuccessMessage: message, HTTPStatusCode: http.StatusOK})
}

func (srv *Server) setUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error retrieving id from URL: %s", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("UNEXPECTED ERROR"))
		return
	}
	//unmarshaling request to domain.User
	var roleReq setUserRoleRequest
	binder := bind.JSONBinder{}
	err = binder.Bind(&roleReq, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})
	switch status := bindingErrorToHTTPStatus(err); status {
	case http.StatusOK:
	case http.StatusInternalServerError:
		log.Printf("unexpected error binding setUseRoleRequest: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unexpected error parsing request body"))
		return
	default:
		exception.WriteErrorAsJSON(w, status, err)
		return
	}

	err = srv.app.SetUserRole(r.Context(), int32(id), roleReq.Role)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourseNotFound):
			exception.WriteErrorAsJSON(w, http.StatusNotFound, err)
			return
		}
		log.Println("error updating user role", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	writeSuccessOK(w, "ok")
}

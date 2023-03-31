package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yawnak/foodadvisor/pkg/bind"
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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("error retrieving id from URL: %s", err)
		writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("UNEXPECTED ERROR"))
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
		writeErrorAsJSON(w, status, errors.New("unexpected error parsing request body"))
		return
	default:
		writeErrorAsJSON(w, status, err)
		return
	}

	err = srv.app.SetUserRole(r.Context(), int32(id), roleReq.Role)
	if err != nil {
		log.Println("error updating user role", err)
		writeErrorAsJSON(w, http.StatusInternalServerError, fmt.Errorf("unknown error"))
		return
	}
	writeSuccessOK(w, "ok")
}

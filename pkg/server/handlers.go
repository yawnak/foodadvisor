package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type key int

const (
	keyUserId key = iota
	keyRole
)

type setUserRoleRequest struct {
	Role string `json:"role"`
}

//TODO: use bind package to parse request and validate request data

func (srv *Server) setUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("error retrieving id from URL: %s", err)
		writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("UNEXPECTED ERROR"))
	}

	//unmarshaling request to domain.User
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) //set maximum size of request body to 1mb
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() //will not allow sending unknown fields in request
	var roleReq setUserRoleRequest
	err = dec.Decode(&roleReq)
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

	err = srv.app.SetUserRole(r.Context(), int32(id), roleReq.Role)
	if err != nil {
		log.Println("error updating user role", err)
		writeErrorAsJSON(w, http.StatusInternalServerError, fmt.Errorf("unknown error"))
	}
}

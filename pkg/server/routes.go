package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (srv *Server) initAPIRoutes() {
	api := srv.router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/signup", validateContentJSON(srv.signup)).Methods("POST")
	api.HandleFunc("/login", validateContentJSON(srv.login)).Methods("POST")

	api.HandleFunc("/profile", srv.auth(func(w http.ResponseWriter, r *http.Request) {
		id, ok := retrieveUserId(r.Context())
		if !ok {
			writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("ERROR: not authorized"))
			log.Fatalln("srv.auth didn't work. no userid in context")
			return
		}
		fmt.Fprintf(w, "hello user number: %d", id)
	})).Methods("GET")
}

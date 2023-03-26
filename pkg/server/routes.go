package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/yawnak/foodadvisor/pkg/domain"
)

func (srv *Server) initAPIRoutes() {
	api := srv.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/signup", validateContentJSON(srv.signup)).Methods("POST")
	api.HandleFunc("/login", validateContentJSON(srv.login)).Methods("POST")

	api.HandleFunc("/profile", srv.authTokenToContext(func(w http.ResponseWriter, r *http.Request) {
		id, ok := retrieveUserId(r.Context())
		if !ok {
			writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("ERROR: not authorized"))
			log.Fatalln("srv.auth didn't work. no userid in context")
			return
		}
		fmt.Fprintf(w, "hello user number: %d", id)
	})).Methods("GET")

	roles := api.PathPrefix("/roles").Subrouter()
	roles.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		role, ok := retrieveRole(r.Context())
		log.Println("ROLES")
		if !ok {
			writeErrorAsJSON(w, http.StatusInternalServerError, errors.New("not authorized"))
			log.Fatalln("didn't work. no role in context")
			return
		}
		fmt.Fprintf(w, "hello user with role: %s", role.Name)
	}).Methods("GET")
	roles.Use(srv.authenticate)
	roles.Use(confirmPermissions(domain.PermEditRoles))
}

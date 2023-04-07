package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/yawnak/foodadvisor/internal/domain"
)

func (srv *Server) initAPIRoutes() {
	// /api subrouter
	api := srv.router.PathPrefix("/api").Subrouter()

	//sign up endpoints
	api.HandleFunc("/signup", validateContentJSON(srv.signup)).Methods("POST")
	api.HandleFunc("/login", validateContentJSON(srv.login)).Methods("POST")

	// /api/users/ routes
	//users := api.PathPrefix("/users").Subrouter()
	usersAuth := api.PathPrefix("/users").Subrouter()

	usersAuth.Use(srv.authenticate)

	usersAuthAndPerms := usersAuth.PathPrefix("").Subrouter()
	usersAuthAndPerms.Use(confirmPermissions(domain.PermEditUserRole))
	usersAuthAndPerms.HandleFunc("/{id:[0-9]+}/role", srv.setUserRole).Methods("POST")

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

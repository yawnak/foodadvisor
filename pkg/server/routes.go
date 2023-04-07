package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
	"github.com/yawnak/foodadvisor/pkg/server/middleware"
)

func (srv *Server) initAPIRoutes() {
	// /api subrouter
	api := srv.router.PathPrefix("/api").Subrouter()

	//sign up endpoints
	api.HandleFunc("/signup", middleware.ValidateContentJSON(srv.signup)).Methods("POST")
	api.HandleFunc("/login", middleware.ValidateContentJSON(srv.login)).Methods("POST")

	// /api/users/ routes
	//users := api.PathPrefix("/users").Subrouter()
	usersAuth := api.PathPrefix("/users").Subrouter()

	usersAuth.Use(srv.authenticate)

	usersAuthAndPerms := usersAuth.PathPrefix("").Subrouter()
	usersAuthAndPerms.Use(middleware.ConfirmPermissions(domain.PermEditUserRole))
	usersAuthAndPerms.HandleFunc("/{id:[0-9]+}/role", srv.setUserRole).Methods("POST")

	roles := api.PathPrefix("/roles").Subrouter()
	roles.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		role, ok := middleware.RetrieveRole(r.Context())
		log.Println("ROLES")
		if !ok {
			exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("not authorized"))
			log.Fatalln("didn't work. no role in context")
			return
		}
		fmt.Fprintf(w, "hello user with role: %s", role.Name)
	}).Methods("GET")
	roles.Use(srv.authenticate)
	roles.Use(middleware.ConfirmPermissions(domain.PermEditRoles))
}

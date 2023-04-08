package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
	"github.com/yawnak/foodadvisor/pkg/server/middleware"
)

func (srv *Server) initAPIRoutes() http.Handler {
	// /api subrouter
	r := chi.NewRouter()
	r.Use(chimiddleware.AllowContentType("application/json"))

	//sign up endpoints
	r.Post("/signup", srv.signup)
	r.Post("/login", srv.login)

	// // /api/users/ routes
	r.With(srv.authenticate).Mount("/users", srv.initUserRoutes())

	return r
}

func (srv *Server) initUserRoutes() http.Handler {
	r := chi.NewRouter()
	// /{id}
	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		//roles routes
		r.Route("/role", func(r chi.Router) {
			r.Use(middleware.ConfirmPermissions(domain.PermEditUserRole))
			// /{id}}/role get user role
			r.Get("/role", func(w http.ResponseWriter, r *http.Request) {
				role, ok := domain.RoleFromContext(r.Context())
				log.Println("ROLES")
				if !ok {
					exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("not authorized"))
					log.Fatalln("didn't work. no role in context")
					return
				}
				fmt.Fprintf(w, "hello user with role: %s", role.Name)
			})
			// /{id}/role set user role
			r.Post("/role", srv.setUserRole)
		})
	})
	return r
}

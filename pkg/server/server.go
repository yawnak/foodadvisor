package server

import (
	"net/http"

	"github.com/asstronom/foodadvisor/pkg/domain"
	"github.com/gorilla/mux"
)

type Server struct {
	app    domain.Advisor
	router *mux.Router
}

func NewServer(app domain.Advisor) (*Server, error) {
	srv := Server{
		app:    app,
		router: mux.NewRouter(),
	}

	srv.initAPIRoutes()

	return &srv, nil
}

func (srv *Server) ListenAndServe(port string) error {
	return http.ListenAndServe(":"+port, srv.router)
}

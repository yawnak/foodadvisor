package server

func (srv *Server) initAPIRoutes() {
	r := srv.router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/signup", validateContentJSON(srv.signup)).Methods("POST")
	r.HandleFunc("/login", validateContentJSON(srv.login)).Methods("POST")
}

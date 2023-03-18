package server

func (srv *Server) initAPIRoutes() {
	api := srv.router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/signup", validateContentJSON(srv.signup)).Methods("POST")
	api.HandleFunc("/login", validateContentJSON(srv.login)).Methods("POST")
	
}

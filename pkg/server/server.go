package server

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/yawnak/foodadvisor/internal/domain"
)

type Server struct {
	app      domain.Advisor
	router   *chi.Mux
	validate *validator.Validate
}

func NewServer(app domain.Advisor) (*Server, error) {
	srv := Server{
		app: app,
		//router: chi.NewRouter(),
	}

	rr := chi.NewRouter()

	//Add the cors middleware with the required headers
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
	//rr.Use(logResponseBody)
	rr.Use(c.Handler)

	rr.Mount("/api", srv.initAPIRoutes())

	srv.router = rr

	srv.initValidator()
	return &srv, nil
}

func (srv *Server) initValidator() error {
	srv.validate = validator.New()
	srv.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})

	srv.validate.RegisterAlias("usernamemax", "lte=30")
	srv.validate.RegisterAlias("passwordmax", "lte=72")

	srv.validate.RegisterValidation("alphanumwithunderscore", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		rg := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return rg.MatchString(val)
	})
	return nil
}

func (srv *Server) ListenAndServe(port string) error {
	return http.ListenAndServe(":"+port, srv.router)
}

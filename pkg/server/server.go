package server

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/yawnak/foodadvisor/pkg/domain"
)

type Server struct {
	app      domain.Advisor
	router   *mux.Router
	validate *validator.Validate
}

func NewServer(app domain.Advisor) (*Server, error) {
	srv := Server{
		app:    app,
		router: mux.NewRouter(),
	}

	srv.initAPIRoutes()
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

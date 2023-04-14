package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yawnak/foodadvisor/internal/domain"
	"github.com/yawnak/foodadvisor/pkg/bind"
	"github.com/yawnak/foodadvisor/pkg/server/exception"
)

type Successer interface {
	Success() string
	Status() int
}

func writeSuccess(w http.ResponseWriter, success Successer) {
	w.WriteHeader(success.Status())
	_ = json.NewEncoder(w).Encode(success)
}

func writeSuccessOK(w http.ResponseWriter, message string) {
	writeSuccess(w, responseSuccess{SuccessMessage: message, HTTPStatusCode: http.StatusOK})
}

func (srv *Server) setUserRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("error retrieving id from URL: %s", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, errors.New("UNEXPECTED ERROR"))
		return
	}
	//unmarshaling request to domain.User
	var roleReq setUserRoleRequest
	binder := bind.JSONBinder{}
	err = binder.Bind(&roleReq, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})
	switch status := bindingErrorToHTTPStatus(err); status {
	case http.StatusOK:
	case http.StatusInternalServerError:
		log.Printf("unexpected error binding setUseRoleRequest: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unexpected error parsing request body"))
		return
	default:
		exception.WriteErrorAsJSON(w, status, err)
		return
	}

	err = srv.app.SetUserRole(r.Context(), int32(id), roleReq.Role)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourseNotFound):
			exception.WriteErrorAsJSON(w, http.StatusNotFound, err)
			return
		}
		log.Println("error updating user role", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	writeSuccessOK(w, "ok")
}

func (srv *Server) createMeal(w http.ResponseWriter, r *http.Request) {
	var mealreq requestCreateFood
	binder := bind.JSONBinder{}
	err := binder.Bind(&mealreq, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})

	switch status := bindingErrorToHTTPStatus(err); status {
	case http.StatusOK:
	case http.StatusInternalServerError:
		log.Printf("unexpected error binding setUseRoleRequest: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unexpected error parsing request body"))
		return
	default:
		exception.WriteErrorAsJSON(w, status, err)
		return
	}

	//struct validation
	err = srv.validate.Struct(mealreq)
	if err != nil {
		exception.WriteErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	food := domain.Food(mealreq)
	foodid, err := srv.app.CreateFood(r.Context(), &food)
	if err != nil {
		log.Println("error creating food:", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	writeSuccess(w, responseCreateMeal{
		responseSuccess: responseSuccess{
			HTTPStatusCode: http.StatusCreated,
			SuccessMessage: "ok",
		},
		MealId: foodid,
	})
}

func (srv *Server) updateUserEaten(w http.ResponseWriter, r *http.Request) {
	userid, ok := userIdFromContext(r.Context())
	if !ok {
		log.Println("updateUserEaten error: no userid in request context")
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
	}
	temp, err := strconv.Atoi(chi.URLParamFromCtx(r.Context(), "foodid"))
	if err != nil {
		exception.WriteErrorAsJSON(w, http.StatusBadRequest, fmt.Errorf("bad foodid"))
		return
	}
	foodid := int32(temp)

	err = srv.app.UpdateUserEaten(r.Context(), userid, foodid, nil)
	if err != nil {
		log.Println("updateUserEaten error:", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}

	writeSuccessOK(w, "ok")
}

func (srv *Server) basicAdvise(w http.ResponseWriter, r *http.Request) {
	id, ok := userIdFromContext(r.Context())
	if !ok {
		log.Println("basicAdvise: no userid in context")
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}

	meals, err := srv.app.BasicAdvise(r.Context(), id, 10, 0)
	if err != nil {
		log.Println("basicAdvise: error getting basic advise")
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	writeSuccess(w, responseBasicAdvice{
		responseSuccess: responseSuccess{
			SuccessMessage: "ok",
			HTTPStatusCode: http.StatusOK,
		},
		Meals: meals,
	})
}

func (srv *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userid, ok := userIdFromContext(r.Context())
	if !ok {
		log.Println("getUser no userid in context")
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	user, err := srv.app.GetUserById(r.Context(), userid)
	if err != nil {
		log.Printf("error getUser GetUserById: %s\n", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, err)
		return
	}
	_ = json.NewEncoder(w).Encode(responseGetUserById(*user))
}

func (srv *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	userid, ok := userIdFromContext(r.Context())
	if !ok {
		log.Println("updateUser: no userid in context")
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}
	var updReq requestUpdateUser
	binder := bind.JSONBinder{}
	err := binder.Bind(&updReq, w, r.Body, &bind.Options{
		MaxBytes:              1 << 20,
		DisallowUnknownFields: true,
	})
	switch status := bindingErrorToHTTPStatus(err); status {
	case http.StatusOK:
	case http.StatusInternalServerError:
		log.Printf("unexpected error binding setUseRoleRequest: %s", err)
		exception.WriteErrorAsJSON(w, status, errors.New("unexpected error parsing request body"))
		return
	default:
		exception.WriteErrorAsJSON(w, status, err)
		return
	}

	err = srv.validate.Struct(updReq)
	if err != nil {
		exception.WriteErrorAsJSON(w, http.StatusBadRequest, err)
		return
	}

	user := domain.User(updReq)

	err = srv.app.UpdateUserById(r.Context(), userid, &user)
	if err != nil {
		log.Printf("updateUser: unknown error updating: %s\n", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, domain.ErrUnknownError)
		return
	}

	writeSuccessOK(w, "ok")
}

func (srv *Server) getMeals(w http.ResponseWriter, r *http.Request) {
	var limit, offset uint
	temp := chi.URLParamFromCtx(r.Context(), "offset")
	if temp == "" {
		offset = 0
	}
	temp = chi.URLParamFromCtx(r.Context(), "limit")
	if temp == "" {
		limit = 20
	}

	meals, err := srv.app.GetMeals(r.Context(), offset, limit)
	if err != nil {
		log.Printf("error getting meals: %s\n", err)
		exception.WriteErrorAsJSON(w, http.StatusInternalServerError, err)
	}

	writeSuccess(w, responseGetMeals{
		responseSuccess: responseSuccess{
			SuccessMessage: "ok",
			HTTPStatusCode: http.StatusOK,
		},
		Meals: meals,
	})
}

package server

import (
	"errors"
	"net/http"

	"github.com/xorcare/pointer"
	"github.com/yawnak/foodadvisor/pkg/bind"
)



// return http.StatusOK if err = nil
// return http.StatusInternalServerError if error is unknown
// other - if error is known (mostly http.StatusBadRequest)
func bindingErrorToHTTPStatus(err error) int {
	if err != nil {
		switch {
		case errors.Is(err, bind.ErrBadFormat):
			return http.StatusBadRequest
		case errors.Is(err, bind.ErrBodyTooLarge):
			return http.StatusBadRequest
		case errors.Is(err, bind.ErrEmptyBody):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrSyntax{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnmarshalType{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnknownField{})):
			return http.StatusBadRequest
		case errors.As(err, pointer.Of(&bind.ErrUnknown{})):
			return http.StatusInternalServerError
		}
	}
	return http.StatusOK
}

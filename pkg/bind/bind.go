package bind

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Binder interface {
	Bind(w http.ResponseWriter, body io.ReadCloser, options Options)
}

type Options struct {
	MaxBytes              int64
	DisallowUnknownFields bool
}

var (
	defaultOptions = Options{
		MaxBytes:              0,
		DisallowUnknownFields: false,
	}
)

type JSONBinder struct {
}

func (b *JSONBinder) Bind(dest any, w http.ResponseWriter, body io.ReadCloser, options Options) error {
	//unmarshaling request to domain.User
	if options.MaxBytes != 0 {
		body = http.MaxBytesReader(w, body, 1<<20) //set maximum size of request body to 1mb
	}
	dec := json.NewDecoder(body)
	if options.DisallowUnknownFields {
		dec.DisallowUnknownFields() //will not allow sending unknown fields in request
	}
	err := dec.Decode(&dest)
	//processing errors
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			return &ErrSyntax{Offset: syntaxError.Offset}
		case errors.Is(err, io.ErrUnexpectedEOF):
			return ErrBadFormat
		case errors.As(err, &unmarshalTypeError):
			return &ErrUnmarshalType{
				Field:  unmarshalTypeError.Field,
				Offset: unmarshalTypeError.Offset,
				Type:   unmarshalTypeError.Type.Name(),
			}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return &ErrUnkownField{Field: fieldName}
		case errors.Is(err, io.EOF):
			return ErrEmptyBody
		case err.Error() == "http: request body too large":
			return ErrBodyTooLarge
		default:
			return &ErrUnknown{Err: err}
		}
	}
	return nil
}

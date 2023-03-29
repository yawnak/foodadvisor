package bind

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xorcare/pointer"
)

type Test struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

type Test1 struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func TestJSONBind(t *testing.T) {
	binder := JSONBinder{}
	type TestData struct {
		dest    any
		data    map[string]interface{}
		options Options
		expErr  error
	}
	data := []TestData{
		{
			dest: &Test{},
			data: map[string]interface{}{
				"foo":  "foo",
				"bar":  "bar",
				"food": "food",
			},
			expErr: nil,
		},
		{
			dest: &Test{},
			data: map[string]interface{}{
				"foo":  "foo",
				"bar":  "bar",
				"food": "food",
			},
			options: Options{DisallowUnknownFields: true},
			expErr:  &ErrUnknownField{Field: "food"},
		},
		{
			dest: &Test1{},
			data: map[string]interface{}{
				"foo":  "foo",
				"bar":  "bar",
				"food": "food",
			},
			options: Options{DisallowUnknownFields: true},
			expErr:  &ErrUnmarshalType{Field: "bar", Offset: 12, Type: "int"},
		},
	}
	for i, v := range data {
		b, err := json.Marshal(v.data)
		if err != nil {
			t.Errorf("test %d. error marshaling map: %s", i, err)
		}
		buff := bytes.NewBuffer(b)
		req, err := http.NewRequest("", "", buff)
		if err != nil {
			t.Errorf("test %d. error making request: %s", i, err)
		}
		rec := &httptest.ResponseRecorder{}
		err = binder.Bind(v.dest, rec, req.Body, &v.options)
		if err != nil {
			switch {
			case v.expErr == nil:
				t.Errorf("test %d. error binding: %v when expected nil\n", i, err)
			case errors.Is(err, v.expErr):
				t.Errorf("test %d. error binding: %v when expected: %v\n", i, err, v.expErr)
			}
		}
		fmt.Println()
		fmt.Printf("test %d PASS.\nbinding error: %v,\nexpected error: %v\n", i, err, v.expErr)
		fmt.Println("res:", v.dest)
	}
}

func TestStuff(t *testing.T) {
	err := error(&ErrUnknownField{Field: "pudge"})

	fmt.Println(errors.As(err, pointer.Of(&ErrUnknownField{})))
}

package clients

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/render"
	"github.com/leonardochaia/vendopunkto/errors"
	"github.com/leonardochaia/vendopunkto/testutils"
)

type dto struct {
	Foo string `json:"string"`
}

func TestHTTPGetWithResult(t *testing.T) {

	expected := dto{
		Foo: "bar",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		testutils.Equals(t, r.URL.String(), "/some/path")

		render.JSON(w, r, expected)
	}))
	defer server.Close()

	client := NewHTTPClient()
	var result dto
	resp, err := client.Get(server.URL+"/some/path", &result)

	testutils.Ok(t, err)
	testutils.Equals(t, expected, result)
	testutils.Equals(t, resp.StatusCode, 200)
}

func TestHTTPGetWithError(t *testing.T) {

	const k = errors.Internal

	server := httptest.NewServer(http.HandlerFunc(
		errors.WrapHandler(
			func(w http.ResponseWriter, r *http.Request) error {
				const op errors.Op = "api.test"
				err := errors.Str("inner")
				return errors.E(op, k, err)
			})))
	defer server.Close()

	client := NewHTTPClient()
	var result dto
	resp, err := client.Get(server.URL+"/some/path", &result)

	testutils.Assert(t, resp.StatusCode >= 400, "expected an error status code")
	testutils.Assert(t, err != nil, "expected an error, got nil")

	e, ok := err.(*errors.Error)
	testutils.Assert(t, ok, "expected an errors.Error")
	testutils.Assert(t, e.Kind == k, "expected the same kind")
}

package rest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteError(t *testing.T) {
	t.Run("unexpected error is encoded as Interlnal Server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		unexpectedError := errors.New("Unexpected error")
		WriteError(w, r, unexpectedError)

		res := w.Result()
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		require.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.JSONEq(t, `{"code": 500, "description": "Internal Server Error"}`, string(body))
	})
	t.Run("HTTPError is encoded in JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		apiError := &HTTPError{
			Code:        http.StatusTeapot,
			Description: "I'm a teapot",
		}
		WriteError(w, r, apiError)

		res := w.Result()
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		require.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.JSONEq(t, `{"code": 418, "description": "I'm a teapot"}`, string(body))
	})
}

func TestAPIHandler(t *testing.T) {
	t.Run("encodes all unexpected errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		var fn APIHandler = func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("Unexpected error")
		}
		fn.ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		require.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.JSONEq(t, `{"code": 500, "description": "Internal Server Error"}`, string(body))
	})
}

func TestMiddlewareHandler(t *testing.T) {
	t.Run("encodes all unexpected errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		var fn MiddlewareHandler = func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error) {
			return w, r, errors.New("Unexpected error")
		}
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.FailNow(t, "Must not be called")
		})
		fn.Middleware(next).ServeHTTP(w, r)

		res := w.Result()
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		require.Equal(t, "application/json", res.Header.Get("Content-Type"))
		assert.JSONEq(t, `{"code": 500, "description": "Internal Server Error"}`, string(body))
	})
	t.Run("can modify request and set default response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		var fn MiddlewareHandler = func(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error) {
			r.Header.Add("X-My-Test-Header", "MODIFIED")
			w.Header().Add("X-Request-ID", "DEFAULT")
			return w, r, nil
		}
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "MODIFIED", r.Header.Get("X-My-Test-Header"))
			w.WriteHeader(http.StatusTeapot)
		})
		fn.Middleware(next).ServeHTTP(w, r)

		res := w.Result()
		require.Equal(t, http.StatusTeapot, res.StatusCode)
		require.Equal(t, res.Header.Get("X-Request-ID"), "DEFAULT")
	})
}

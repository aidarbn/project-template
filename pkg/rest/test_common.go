package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	logs "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// testRequest is used to imitate requests to REST API
type testRequest struct {
	params        map[string]any // `any` is used only for primitive types
	headers       map[string]any
	contextParams map[any]any
	pathParams    map[string]any
	bodyReader    io.Reader
}

// WithPathParam is optional argument for testRequest. Represents path parameter.
func WithPathParam(key string, value any) func(*testRequest) {
	return func(r *testRequest) {
		r.pathParams[key] = value
	}
}

// WithQuery is optional argument for testRequest. Represents query parameter.
func WithQuery(key string, value any) func(*testRequest) {
	return func(r *testRequest) {
		r.params[key] = value
	}
}

// WithHeader is optional argument for testRequest. Represents header parameter.
func WithHeader(key string, value any) func(*testRequest) {
	return func(r *testRequest) {
		r.headers[key] = value
	}
}

// WithContextParam is optional argument for testRequest. Represents context parameter added to *http.Request.
func WithContextParam(key any, value any) func(*testRequest) {
	return func(r *testRequest) {
		r.contextParams[key] = value
	}
}

// WithBody is optional argument for testRequest. Represents body of the request.
func WithBody(body []byte) func(*testRequest) {
	return func(r *testRequest) {
		r.bodyReader = bytes.NewReader(body)
	}
}

// WithJSON is optional argument for testRequest. Represents body of the request in the form of JSON.
func WithJSON(object any) func(*testRequest) {
	return func(r *testRequest) {
		body, err := json.Marshal(object)
		if err != nil {
			logs.Error(err)
		}
		r.bodyReader = bytes.NewReader(body)
		r.headers[ContentType] = ContentTypeJSON
	}
}

// WithBodyReader is optional argument for testRequest. Represents body of the request in any form type.
// Helpful for multipart files sending.
func WithBodyReader(reader io.Reader) func(*testRequest) {
	return func(r *testRequest) {
		r.bodyReader = reader
	}
}

func WithEmptyBody() func(*testRequest) {
	return func(r *testRequest) {
		r.bodyReader = nil
	}
}

// NewTestRequest imitates REST requests with optional variables.
func NewTestRequest(options ...func(*testRequest)) (*httptest.ResponseRecorder, *http.Request) {
	tr := &testRequest{
		params:        make(map[string]any),
		headers:       make(map[string]any),
		contextParams: make(map[any]any),
		pathParams:    make(map[string]any),
		bodyReader:    nil,
	}
	for _, opt := range options {
		opt(tr)
	}
	paramsVals := url.Values{}
	r := httptest.NewRequest(http.MethodGet, "/", tr.bodyReader)
	for key, val := range tr.params {
		paramsVals.Set(key, fmt.Sprintf("%+v", val))
	}
	r.URL.RawQuery = paramsVals.Encode()

	for key, val := range tr.headers {
		r.Header.Set(key, fmt.Sprintf("%+v", val))
	}

	for key, val := range tr.contextParams {
		ctx := context.WithValue(r.Context(), key, val)
		r = r.WithContext(ctx)
	}

	for key, val := range tr.contextParams {
		ctx := context.WithValue(r.Context(), key, val)
		r = r.WithContext(ctx)
	}

	chiCtx := chi.NewRouteContext()
	for key, val := range tr.pathParams {
		chiCtx.URLParams.Add(key, fmt.Sprintf("%+v", val))
	}
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

	w := httptest.NewRecorder()

	return w, r
}

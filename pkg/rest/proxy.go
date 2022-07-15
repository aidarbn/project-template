package rest

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy is proxying service.
type Proxy struct {
	*chi.Mux

	defaultRequestData proxyRequestData // data that can modify request
}

func NewProxy(router *chi.Mux, defaultOpts ...Option) *Proxy {
	proxy := &Proxy{
		Mux: router,
	}
	proxy.AddDefaultOptions(defaultOpts...)
	return proxy
}

// AddDefaultOptions adds default options.
func (p *Proxy) AddDefaultOptions(opts ...Option) {
	for _, opt := range opts {
		opt.Add(&p.defaultRequestData)
	}
}

// ProxyRequest performs request.
func (p *Proxy) ProxyRequest(method, origPath string, opts ...Option) {
	proxyData := proxyRequestData{}
	for _, opt := range opts {
		opt.Add(&proxyData)
	}
	p.Mux.Method(method, origPath,
		&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				for _, director := range p.defaultRequestData.directors {
					director(req)
				}
				for _, director := range proxyData.directors {
					director(req)
				}
			},
			ModifyResponse: func(resp *http.Response) error {
				for _, modifier := range p.defaultRequestData.modifiers {
					if err := modifier(resp); err != nil {
						return err
					}
				}
				for _, modifier := range proxyData.modifiers {
					if err := modifier(resp); err != nil {
						return err
					}
				}
				return nil
			},
		},
	)
}

// SetBaseURL sets url base for current path.
func SetBaseURL(baseURL *url.URL) func(r *http.Request) {
	return func(r *http.Request) {
		r.Host = baseURL.Host
		r.URL.Scheme = baseURL.Scheme
		r.URL.Host = baseURL.Host
	}
}

// TrimPathPrefix deleted prefix from current path.
func TrimPathPrefix(prefix string) func(r *http.Request) {
	return func(r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
	}
}

// proxyRequestData is a data needed for proxy request.
type proxyRequestData struct {
	directors []Director
	modifiers []ResponseModifier
}

// Option is a function that should be performed at proxy request.
type Option interface {
	// Add adds option to the proxy service.
	Add(*proxyRequestData)
}

// Director modifies proxy request before it proxies.
type Director func(req *http.Request)

func (d Director) Add(proxyData *proxyRequestData) {
	proxyData.directors = append(proxyData.directors, d)
}

// ResponseModifier modifies response of the request.
type ResponseModifier func(resp *http.Response) error

func (r ResponseModifier) Add(proxyData *proxyRequestData) {
	proxyData.modifiers = append(proxyData.modifiers, r)
}

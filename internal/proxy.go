package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// TargetRoutes
type Route struct {
	// optional a name for route for the sake of conviency
	Name string
	// available servers for loadbalancing
	// should be at minimum 1 server
	Pool *Pool

	// the route path
	// eg. /health
	Path string
}

// pool or hosts for loadbalancing
type Backend struct {
	Url       *url.URL
	Proxy     *httputil.ReverseProxy
	IsHealthy bool
	Name      string //optional for debugging
}
type Pool struct {
	Backends []*Backend
	Next     int32
	//mu      sync.Mutex // to synchronize access , in the context of conccurent requests
}

func NewPool(hosts []string) (*Pool, error) {
	var backends []*Backend = make([]*Backend, len(hosts))
	for idx, host := range hosts {
		backend := Backend{
			Name: fmt.Sprintf("backend id[ %d ] Host[ %s]", idx, host)}
		url, err := url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf(" backend url wrong: %s", host)
		}
		backend.Url = url
		//TODO: add way to ping the backend to check health
		backend.IsHealthy = true
		backend.Proxy = httputil.NewSingleHostReverseProxy(backend.Url)
		backends[idx] = &backend
	}
	pool := Pool{
		Backends: backends,
		Next:     0,
	}
	return &pool, nil

}

func NewRoute(config RouteConfig) (*Route, error) {
	pool, err := NewPool(config.Backends)
	if err != nil {
		return nil, errors.New("error initializing the pool")
	}
	route := Route{
		Path: config.Path,
		Pool: pool,
	}
	return &route, nil
}

// round robin load balancing
func (p *Pool) NextProxy() *httputil.ReverseProxy {

	proxy := p.Backends[p.Next].Proxy
	//TODO: add checking the ishealthy condition
	p.Next = atomic.AddInt32(&p.Next, int32(1)) % int32(len(p.Backends))
	return proxy
}

// wrapper for pool.NextProxy()
func (r *Route) NextProxy() *httputil.ReverseProxy {
	return r.Pool.NextProxy()
}

// handle the req by proxying it to some backend from the pool
func (route *Route) HandleReqOverProxy() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request ack .")
		fmt.Println("path:", r.Host, r.RemoteAddr, r.RequestURI)
		proxy := route.NextProxy()
		proxy.ServeHTTP(w, r)

	}
}

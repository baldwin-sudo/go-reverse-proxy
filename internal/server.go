package proxy

import (
	"fmt"
	"net/http"
)

// the proxy server
type Server struct {
	Port   string
	Routes []*Route
	Config Config
	Mux    *http.ServeMux
	http.Server
}

// Initialize and bind routes to the mux
func (s *Server) HandleAllRoutes() {
	s.Routes = make([]*Route, 0)
	s.Mux = http.NewServeMux()

	for _, routeConfig := range s.Config.RoutesConfig {
		route, err := NewRoute(routeConfig)
		if err != nil {
			fmt.Println("error initializing route:", routeConfig.Path)
			continue
		}
		s.Routes = append(s.Routes, route)
		s.Mux.Handle(route.Path, route.HandleReqOverProxy())
	}

	// Attach the mux to the embedded http.Server
	s.Server = http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Mux,
	}
}

// Start listening
func (s *Server) Start() error {
	fmt.Println("Starting proxy on port", s.Port)
	return s.ListenAndServe()
}

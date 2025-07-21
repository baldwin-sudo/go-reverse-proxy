package proxy

import (
	"fmt"
	"net/http"
)

// the proxy server
type Server struct {
	http.Server
	Port   string
	Routes []*Route
	Config Config
}

func (server *Server) HandleAllRoutes() {
	for _, routeConfig := range server.Config.RoutesConfig {

		route, err := NewRoute(routeConfig)
		if err != nil {
			fmt.Println("error initializing routes")
			return
		}
		http.Handle(route.Path, route.HandleReqOverProxy())

	}
}

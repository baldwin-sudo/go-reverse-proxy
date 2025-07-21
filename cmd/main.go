package main

import (
	// Replace the import path below with the actual path to your internal/proxy package

	"fmt"
	"net/http"

	"github.com/baldwin-sudo/go-reverse-proxy/internal/proxy"
)

func main() {
	//example
	config := proxy.Config{
		ServerConfig: proxy.ServerConfig{
			Port:    "8080",
			Address: "localhost",
		},
		RoutesConfig: []proxy.RouteConfig{
			proxy.RouteConfig{
				Path: "/health1",
				Backends: []string{
					"http://localhost:8081",
					"http://localhost:8082"},
			},
			proxy.RouteConfig{
				Path: "/health2",
				Backends: []string{
					"http://localhost:8083",
					"http://localhost:8084"},
			},
		},
	}
	server := proxy.Server{
		Port: config.ServerConfig.Port,
	}
	for _, routeCfg := range config.RoutesConfig {
		route, err := proxy.NewRoute(routeCfg)
		if err != nil {
			fmt.Println("error creating route :", routeCfg)
			return
		}
		http.Handle(route.Path, route.HandleReqOverProxy())
	}
	http.ListenAndServe(":"+server.Port, nil)
}

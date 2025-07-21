package main

import (
	// Replace the import path below with the actual path to your internal/proxy package

	"fmt"
	"net/http"

	proxy "github.com/baldwin-sudo/go-reverse-proxy/internal"
)

func main() {
	proxy.LoadConfig("config_example.json")
	server := proxy.Server{
		Config: proxy.Cfg,
		Port:   proxy.Cfg.ServerConfig.Port,
	}
	fmt.Println("starting the reverse proxy ...")
	server.HandleAllRoutes()
	server.Mux.HandleFunc("/proxy_health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	if err := server.Start(); err != nil {
		fmt.Println("error serving the proxy ...")
		return
	}
}

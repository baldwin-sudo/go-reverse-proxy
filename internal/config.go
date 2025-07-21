package proxy

import (
	"encoding/json"
	"os"
)

var Cfg Config

type Config struct {
	ServerConfig ServerConfig  `json:"server"`
	RoutesConfig []RouteConfig `json:"routes"`
}
type ServerConfig struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}
type RouteConfig struct {
	Path     string   `json:"path"`
	Backends []string `json:"backends"`
}

func LoadConfig(configFile string) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic("error reding the file")
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic("error parsing json")
	}
	Cfg = config
}

package proxy

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

}

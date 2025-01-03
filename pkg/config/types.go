package config

type ReverseproxyConfig struct {
	ListenPort  string   `yaml:"listenPort,omitempty"`
	HttpsRoutes []*Route `yaml:"https,omitempty"`
	HttpRoutes  []*Route `yaml:"http,omitempty"`
}

type Route struct {
	Host    string `yaml:"host"`
	Backend string `yaml:"backend"`
	Path    string `yaml:"pathPrefix"`
	Port    int    `yaml:"port"`
}

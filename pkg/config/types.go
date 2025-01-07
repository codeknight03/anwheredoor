package config

type ReverseproxyConfig struct {
	ListenPort string            `yaml:"listenPort,omitempty"`
	Routes     []*Route          `yaml:"routes,omitempty"`
	TLSConfig  *TLSConfiguration `yaml:"tlsConfig,omitempty"`
}

type TLSConfiguration struct {
	// Default certificate to use if no host match is found
	DefaultCertificate *CertificateConfig `yaml:"defaultCertificate,omitempty"`
	// Map of host to certificate configurations
	HostCertificates map[string]*CertificateConfig `yaml:"hostCertificates,omitempty"`
	// Optional: Global TLS settings
	MinVersion string `yaml:"minVersion,omitempty"` // e.g., "1.2", "1.3"
	MaxVersion string `yaml:"maxVersion,omitempty"`
	// Optional: Enable automatic HTTP to HTTPS redirection

	EnableHTTPRedirect bool `yaml:"enableHTTPRedirect,omitempty"`
}

type CertificateConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
	// Optional: Include intermediate certificates
	ChainFile string `yaml:"chainFile,omitempty"`
}

type Route struct {
	Host      string `yaml:"host"`
	Backend   string `yaml:"backend"`
	Path      string `yaml:"pathPrefix"`
	Port      int    `yaml:"port"`
	EnableTLS bool   `yaml:"enableTLS,omitempty"`
	// Optional: Force redirection to https for this route.
	// This overwrites the EnableHTTPRedirect.
	EnableHTTPRedirect bool `yaml:"enableHTTPRedirect,omitempty"`
}

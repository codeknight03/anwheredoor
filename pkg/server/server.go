package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/codeknight03/anywheredoor/pkg/config"
)

type ReverseProxy struct {
	routes map[string]http.Handler
}

func NewReverseProxy(config *config.ReverseproxyConfig) *ReverseProxy {
	routes := make(map[string]http.Handler)

	for _, route := range config.HttpRoutes {
		backendUrl, err := url.Parse("http://" + route.Backend + ":" + string(route.Port))
		if err != nil {
			fmt.Printf("Dropping %v because due to malformed url.", route)
		}

		routes[route.Host+route.Path] = &BackendHandler{
			Target: backendUrl,
		}
	}

	//handle https routes separately until ssl termination is implemented
	//TODO: Move HTTP and HTTPS routes together when SSL termination is implemented

	for _, route := range config.HttpsRoutes {
		backendUrl, err := url.Parse("https://" + route.Backend + ":" + string(route.Port))
		if err != nil {
			fmt.Printf("Dropping %v because due to malformed url.", route)
		}

		routes[route.Host+route.Path] = &BackendHandler{
			Target: backendUrl,
		}
	}

	return &ReverseProxy{routes: routes}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := rp.routes[r.Host+r.URL.Path]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}

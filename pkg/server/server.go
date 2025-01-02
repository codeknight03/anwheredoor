package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/codeknight03/anywheredoor/pkg/config"
)

type ReverseProxy struct {
	routes map[string][]pathHandlerPair
}

type pathHandlerPair struct {
	path    string
	handler http.Handler
}

func NewReverseProxy(config *config.ReverseproxyConfig) *ReverseProxy {

	// Rationale: The grouping of routes needs to be done based on the Host
	//            so that eventually we can design an efficient way to run
	//            more than one listener per Host.
	// Decision:  Map from Host to path and traverse paths in the order of
	//            their length to prefer the most specific path.

	routes := make(map[string][]pathHandlerPair)

	for _, route := range config.HttpRoutes {
		backendUrl, err := url.Parse("http://" + route.Backend + ":" + fmt.Sprint(route.Port))
		if err != nil {
			fmt.Printf("Dropping %v because due to malformed url.", route)
		}

		routes[route.Host] = addToPathHandlerMap(routes[route.Host], route.Path, &BackendHandler{
			Target: backendUrl,
		})
	}

	//handle https routes separately until ssl termination is implemented
	//TODO: Move HTTP and HTTPS routes together when SSL termination is implemented

	for _, route := range config.HttpsRoutes {
		backendUrl, err := url.Parse("https://" + route.Backend + ":" + fmt.Sprint(route.Port))
		if err != nil {
			fmt.Printf("Dropping %v because due to malformed url.", route)
		}

		routes[route.Host] = addToPathHandlerMap(routes[route.Host], route.Path, &BackendHandler{
			Target: backendUrl,
		})
	}

	return &ReverseProxy{routes: routes}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathHandlers, ok := rp.routes[r.Host]
	if !ok {
		fmt.Printf("Host Not Found :%s", r.Host)
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	var handler http.Handler
	for _, pathHandler := range pathHandlers {
		if strings.Contains(r.URL.Path, pathHandler.path) {
			handler = pathHandler.handler
			break
		}
	}

	if handler == nil {
		fmt.Printf("Host and Path Combination not found %s:%s", r.Host, r.URL.Path)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}

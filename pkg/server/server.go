package server

import (
	"fmt"
	"log/slog"
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

func makeRoutesFromConfig(config *config.ReverseproxyConfig) map[string][]pathHandlerPair {
	// Rationale: The grouping of routes needs to be done based on the Host
	//            so that eventually we can design an efficient way to run
	//            more than one listener per Host.
	// Decision:  Map from Host to path and traverse paths in the order of
	//            their length to prefer the most specific path.

	routes := make(map[string][]pathHandlerPair)
	for _, route := range config.HttpRoutes {
		backendUrl, err := url.Parse("http://" + route.Backend + ":" + fmt.Sprint(route.Port))
		if err != nil {
			slog.Warn("Dropping route because due to malformed url.", "route", route, "error", err)
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
			slog.Warn("Dropping route because due to malformed url.", "route", route, "error", err)
		}

		routes[route.Host] = addToPathHandlerMap(routes[route.Host], route.Path, &BackendHandler{
			Target: backendUrl,
		})
	}

	return routes
}
func NewReverseProxy(config *config.ReverseproxyConfig) *ReverseProxy {

	routes := makeRoutesFromConfig(config)

	return &ReverseProxy{routes: routes}
}

func (rp *ReverseProxy) UpdateConfig(config *config.ReverseproxyConfig) *ReverseProxy {

	rp.routes = makeRoutesFromConfig(config)

	return rp
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathHandlers, ok := rp.routes[r.Host]
	if !ok {
		slog.Debug("Host Not Found", "host", r.Host)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var handler http.Handler
	for _, pathHandler := range pathHandlers {
		if strings.Contains(r.URL.Path, pathHandler.path) {
			handler = pathHandler.handler
			break
		}
	}

	if handler == nil {
		slog.Debug("Host and Path Combination not found", "host", r.Host, "path", r.URL.Path)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}

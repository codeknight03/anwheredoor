package server

import "net/http"

// Create separate handler types for TLS check and HTTPS redirect
type TLSCheckHandler struct {
	next http.Handler
}

func (h *TLSCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		http.Error(w, "TLS Required", http.StatusForbidden)
		return
	}
	h.next.ServeHTTP(w, r)
}

type HTTPSRedirectHandler struct {
	next http.Handler
	host string
}

func (h *HTTPSRedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		target := "https://" + h.host + r.RequestURI
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
		return
	}
	h.next.ServeHTTP(w, r)
}

// Modify the wrapper functions to use these handlers
func wrapWithTLSCheck(handler http.Handler) http.Handler {
	return &TLSCheckHandler{next: handler}
}

func wrapWithHTTPSRedirect(handler http.Handler, host string) http.Handler {
	return &HTTPSRedirectHandler{
		next: handler,
		host: host,
	}
}

package server

import (
	"fmt"
	"net/http"
	"net/url"
)

type BackendHandler struct {
	Target *url.URL
}

func (bh *BackendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest(r.Method, bh.Target.String(), r.Body)
	if err != nil {
		fmt.Printf("Error creating request for %s: %s", bh.Target, err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	//copy the headers over to the new request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error in the response for %s:%s", bh.Target, err)
		http.Error(w, "Error forwarding the request", http.StatusBadGateway)
	}

	defer resp.Body.Close()

	//copy the headers back to original response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

}

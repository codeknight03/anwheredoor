package server

import (
	"fmt"
	"io"
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

	//must convert to bytes before sending
	//TODO: Look into this wasted conversion into bytes for writing to a stream again.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading the response for %s: %s\n", bh.Target, err)
		http.Error(w, "Error reading the response body", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the headers to the original response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the status code
	w.WriteHeader(resp.StatusCode)

	// Write the body
	_, writeErr := w.Write(body)
	if writeErr != nil {
		fmt.Printf("Error writing the response body for %s: %s\n", bh.Target, writeErr)
	}

}

package server

import (
	"io"
	"net/http"
	"net/url"

	"log/slog"
)

type BackendHandler struct {
	Target *url.URL
}

func (bh *BackendHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest(r.Method, bh.Target.String(), r.Body)
	if err != nil {
		slog.Debug("Error creating request.", "target", bh.Target, "error", err)
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
		slog.Debug("Error in the response.", "target", bh.Target, "error", err)
		http.Error(w, "Error forwarding the request", http.StatusBadGateway)
	}

	//must convert to bytes before sending
	//TODO: Look into this wasted conversion into bytes for writing to a stream again.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Debug("Error reading the response", "target", bh.Target, "error", err)
		http.Error(w, "Error reading the response body", http.StatusBadGateway)
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
		slog.Debug("Error writing the response body.", "target", bh.Target, "error", writeErr)
		http.Error(w, "Error writing the response body", http.StatusBadGateway)
	}

}

package main

import (
	"fmt"
	"io"
	"net/http"
)

type ConnectionHandler struct{}

func (ConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Copy headers
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// Create a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// copy headers
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	// Adds "proxy: conduit" header
	w.Header().Set("proxy", "conduit")

	// copy status code
	w.WriteHeader(resp.StatusCode)

	// copy body
	io.Copy(w, resp.Body)
}

func main() {
	// Creates a connection handler
	handleTunnelRequest := ConnectionHandler{}
	// Creates a http server to act as proxy
	err := http.ListenAndServe(":8080", handleTunnelRequest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server started on port 8080")
}

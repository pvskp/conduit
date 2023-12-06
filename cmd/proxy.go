/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var (
	port int
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts a HTTP proxy server",
	Long: `Starts a HTTP proxy server on the specified port. The proxy forwards all incoming HTTP requests to the target hosts defined in the configuration.

Usage example:
$ conduit proxy -p 8080

This command will start the proxy server listening on port 8080.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Initiating proxy mode...")
		// Creates a connection handler
		handleTunnelRequest := proxyHandler{}

		// Creates a http server to act as proxy
		log.Println("Listening on port", port, "...")
		portString := fmt.Sprintf(":%d", port)
		err := http.ListenAndServe(portString, handleTunnelRequest)
		if err != nil {
			panic(err)
		}
	},
}

type proxyHandler struct{}

func (proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	if !strings.HasPrefix(r.URL.String(), "http://") {
		protocol := "http://"
		if r.TLS != nil {
			protocol = "https://"
		}
		url = protocol + url
		if protocol == "" {
			log.Println("Using protocol http")
		} else {
			log.Println("Using protocol https") //TODO: HTTPS not implemented yet
		}
	}
	log.Printf("Received new HTTP connection to %s: from %s\n", r.URL.Host, r.RemoteAddr)
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		log.Printf("Error on creating Request object: '%s'\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// Create a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending the request: '%s'\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Request to %s returned status code %d\n", r.Host, resp.StatusCode)
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

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

type LbAlgorithm int

const (
	RoundRobin LbAlgorithm = iota
	LeastConnections
	LeastTime
)

type LoadBalancer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type RRLoadBalancer struct {
	Hosts []url.URL
	Index int
}

func (rr *RRLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New request received. Forwarding to %s\n", rr.Hosts[rr.Index].String())
	reverseProxy := httputil.NewSingleHostReverseProxy(&rr.Hosts[rr.Index])
	rr.NextHost()
	reverseProxy.ServeHTTP(w, r)
}

func (rr *RRLoadBalancer) NextHost() { rr.Index = (rr.Index + 1) % len(rr.Hosts) }

func NewRRLoadBalancer(hosts []url.URL) *RRLoadBalancer {
	return &RRLoadBalancer{
		Hosts: hosts,
	}
}

var (
	LoadBalancerAlgorithm string
	Port                  int
	Hosts                 []string
	algorithmUsage        string = `Possible values:	
roundRobin
leastConnections
leastTime
`
)

// ParseUrl identifies if a given host has a protocol on it's string
// If do not, it adds "http://" before the string
func ParseUrl(host string) string {
	if strings.HasPrefix(host, "http://") {
		return host
	} else if strings.HasPrefix(host, "https://") {
		return host
	}

	return fmt.Sprintf("%s%s", "http://", host)
}

var loadBalancerCmd = &cobra.Command{
	Use:   "loadBalancer",
	Short: "Starts a load balancer",
	Long: `Starts a load balancer that distributes incoming HTTP requests across multiple backend servers. This load balancer supports different algorithms such as Round Robin, Least Connections, and Least Time, providing flexibility in handling traffic.

By default, the load balancer uses the Round Robin algorithm, which distributes requests evenly across all available servers. This is useful for ensuring that no single server bears too much load.

Usage example:
$ conduit loadBalancer -a roundRobin -H http://server1.example.com,http://server2.example.com -p 8080

This command will start the load balancer listening on port 8080 and distribute incoming requests between server1.example.com and server2.example.com using the Round Robin algorithm.`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyFlags(cmd)

		for i, h := range Hosts {
			Hosts[i] = ParseUrl(h)
		}

		urls := parseHosts()
		var lb LoadBalancer
		switch LoadBalancerAlgorithm {
		case "roundRobin":
			log.Println("Starting load balancer with RoundRobin algorithm")
			lb = NewRRLoadBalancer(urls)
			http.HandleFunc("/", lb.ServeHTTP)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); err != nil {
				log.Fatalf("Error starting server: %s", err)
			}

		case "leastConnections":
			//TODO: implement leastConnections
			log.Fatal("LeastConnections algorithm not implemented yet")
		case "leastTime":
			//TODO: implement leastTime
			log.Fatal("LeastTime algorithm not implemented yet")
		}
	},
}

func parseHosts() (urls []url.URL) {
	for _, v := range Hosts {
		u, err := url.Parse(v)
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, *u)
	}
	return urls
}

func verifyFlags(cmd *cobra.Command) {
	if len(Hosts) == 0 {
		log.Fatal("No hosts provided")
		cmd.Usage()
	}

	if LoadBalancerAlgorithm == "" {
		fmt.Println("No loadbalancing algorithm provided")
		fmt.Println(algorithmUsage)
	}

	// verify if the algorithm is RoundRobin, LeastConnections or LeastTime
	if LoadBalancerAlgorithm != "roundRobin" &&
		LoadBalancerAlgorithm != "leastConnections" &&
		LoadBalancerAlgorithm != "leastTime" {
		fmt.Println("Invalid loadbalancing algorithm")
		fmt.Println(algorithmUsage)
	}
}

func init() {
	rootCmd.AddCommand(loadBalancerCmd)
	loadBalancerCmd.Flags().StringSliceVarP(
		&Hosts,
		"hosts",
		"H",
		[]string{},
		"Hosts to load balance",
	)

	loadBalancerCmd.Flags().IntVarP(
		&Port,
		"port",
		"p",
		8080,
		"Port to listen to",
	)

	// list of options for the load balancer algorithm
	loadBalancerCmd.Flags().StringVarP(
		&LoadBalancerAlgorithm,
		"algorithm",
		"a",
		"roundRobin",
		algorithmUsage,
	)
}

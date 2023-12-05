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

// loadBalancerCmd represents the loadBalancer command
var loadBalancerCmd = &cobra.Command{
	Use:   "loadBalancer",
	Short: "loadBalancer starts the load balancer",
	Long:  `loadBalancer starts the load balancer with the specified algorithm and hosts`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyFlags(cmd)

		var lb LoadBalancer

		urls := parseHosts()

		switch LoadBalancerAlgorithm {
		case "roundRobin":
			log.Println("Starting load balancer with RoundRobin algorithm")
			lb = NewRRLoadBalancer(urls)

		case "leastConnections":
			//TODO: implement leastConnections
			log.Fatal("LeastConnections algorithm not implemented yet")
		case "leastTime":
			//TODO: implement leastTime
			log.Fatal("LeastTime algorithm not implemented yet")
		}

		http.Handle("/", lb)
	},
}

func parseHosts() (urls []url.URL) {
	for _, v := range Hosts {
		u, _ := url.Parse(v)
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

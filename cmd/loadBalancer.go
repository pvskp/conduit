/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type RRQueue struct {
	Hosts []string
	Index int
}

var (
	RRQueueInstance       RRQueue = RRQueue{}
	LoadBalancerAlgorithm string
	Client                *http.Client = &http.Client{}
	Port                  int

	algorithmUsage string = `Possible values:	
	roundRobin
	leastConnections
	leastTime
	`
)

// loadBalancerCmd represents the loadBalancer command
var loadBalancerCmd = &cobra.Command{
	Use:   "loadBalancer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		verifyFlags()
		switch LoadBalancerAlgorithm {
		case "roundRobin":
			log.Println("Starting load balancer with RoundRobin algorithm")
			rrhandle := RRHandler{}
			port := fmt.Sprintf(":%d", Port)
			http.ListenAndServe(port, rrhandle)
		case "leastConnections":
			//TODO: implement leastConnections
			log.Fatal("LeastConnections algorithm not implemented yet")
		case "leastTime":
			//TODO: implement leastTime
			log.Fatal("LeastTime algorithm not implemented yet")
		}
	},
}

func verifyFlags() {
	if len(RRQueueInstance.Hosts) == 0 {
		log.Fatal("No hosts provided")
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

type RRHandler struct{}

func (RRHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: implenting parse from a config file to parse possible paths and deny anything different

	parseURL, err := url.Parse(r.URL.String())
	path := parseURL.Path

	// Create the request forwarding it to the next host in the queue
	request, err := http.NewRequest(r.Method, RRQueueInstance.Hosts[RRQueueInstance.Index]+path, r.Body)

	if err != nil {
		log.Println("Error creating request:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := Client.Do(request)
	if err != nil {
		log.Println("Error sending request:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)
}

func (rr *RRQueue) UpdateQueueIndex() { rr.Index = (rr.Index + 1) % len(rr.Hosts) }

// RoundRobin is a function that load balances the hosts equally
func RoundRobin() {}

func init() {
	rootCmd.AddCommand(loadBalancerCmd)
	loadBalancerCmd.Flags().StringSliceVarP(&RRQueueInstance.Hosts, "hosts", "H", []string{}, "Hosts to load balance")
	loadBalancerCmd.Flags().IntVarP(&Port, "port", "p", 8080, "Port to listen to")

	// list of options for the load balancer algorithm
	loadBalancerCmd.Flags().StringVarP(&LoadBalancerAlgorithm, "algorithm", "a", "roundRobin", algorithmUsage)
}

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "conduit",
	Short: "Conduit is a versatile network tool for HTTP routing and load balancing",
	Long: `Conduit is a versatile command-line application designed for HTTP routing and load balancing. It provides various network functionalities including running a simple HTTP proxy and a configurable load balancer.

Example usage:
1. Start a proxy server:
   $ conduit proxy -p 8080

2. Start a load balancer with Round Robin algorithm:
   $ conduit loadBalancer -a roundRobin -H http://server1.example.com,http://server2.example.com -p 8080

Conduit is designed to be lightweight, easy to use, and highly configurable to suit various networking needs.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

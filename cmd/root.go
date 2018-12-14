package cmd

import (
	"fmt"
	"os"

	"github.com/criteo/consul-bench/internal/bench"
	"github.com/spf13/cobra"
)

var benchConfig bench.Config
var bencher *bench.Bench

func init() {
	cobra.OnInitialize(initBench)
	rootCmd.PersistentFlags().StringVarP(&benchConfig.ConsulHost, "consul", "c", "127.0.0.1", "Consul host")
	rootCmd.PersistentFlags().IntVar(&benchConfig.HTTPPort, "http-port", 8500, "Consul http port")
	rootCmd.PersistentFlags().IntVar(&benchConfig.RPCPort, "rpc-port", 8300, "Consul rpc port")
	rootCmd.PersistentFlags().StringVarP(&benchConfig.ACLToken, "token", "t", "", "Consul ACL token")
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A consul benchmarker",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute runs the cli interface
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func initBench() {
	b, err := bench.New(benchConfig)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	bencher = b
}

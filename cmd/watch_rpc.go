package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	watchServiceName string
)

func init() {
	rootCmd.AddCommand(registerCommand)
	registerCommand.Flags().StringVarP(&watchServiceName, "service", "s", "srv", "Service name")
}

var watchRPCCommand = &cobra.Command{
	Use:   "watch-rpc",
	Short: "Register service instances",
	Run: func(cmd *cobra.Command, args []string) {
		err := bencher.RegisterServices(watchServiceName, 1)
		if err != nil {
			log.Println(err)
		}
	},
}

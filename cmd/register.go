package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	registerServiceName string
)

func init() {
	rootCmd.AddCommand(registerCommand)
	registerCommand.Flags().StringVarP(&registerServiceName, "service", "s", "srv", "Service name")
}

var registerCommand = &cobra.Command{
	Use:   "register",
	Short: "Register service instances",
	Run: func(cmd *cobra.Command, args []string) {
		err := bencher.RegisterServices(registerServiceName, 1)
		if err != nil {
			log.Println(err)
		}
	},
}

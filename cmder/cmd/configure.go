package cmd

import (
		"fmt"

		"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "configure for app",
		Long:  `sync config data to etcd configure center`,
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("configure called")
		},
}

func init() {

		configureCmd.PersistentFlags().String("file", "", "action config file path")
		configureCmd.PersistentFlags().String("prefix", "", "action config prefix")

		rootCmd.AddCommand(configureCmd)
}

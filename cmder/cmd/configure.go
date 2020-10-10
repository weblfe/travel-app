package cmd

import (
		"fmt"

		"github.com/spf13/cobra"
)

var (
		file     *string
		prefix   *string
		excludes *[]string
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "configure for app",
		Long:  `sync config data to etcd configure center`,
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("configure called")
				fmt.Println(*file)
				fmt.Println(*prefix)
				fmt.Println(*excludes)
		},
}

func init() {
		file = configureCmd.PersistentFlags().String("file", "", "action config file path")
		prefix = configureCmd.PersistentFlags().String("prefix", "", "action config prefix")
		excludes = configureCmd.PersistentFlags().StringArray("excludes", []string{}, "action config prefix")
		rootCmd.AddCommand(configureCmd)
}

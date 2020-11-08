package cmd

import (
		"github.com/spf13/cobra"
		"github.com/weblfe/travel-app/cmder/kernel"
		"log"
)

var (
		file     string
		prefix   string
		excludes []string
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "configure for app",
		Long:  `sync config data to etcd configure center`,
		Run: func(cmd *cobra.Command, args []string) {
				handler := kernel.InvokerConfigureService(file, prefix, excludes,endpoints,timeout)
				if err := handler.Exec(); err != nil {
						log.Fatal(err)
				}
		},
}

func init() {
		configureCmd.PersistentFlags().StringVar(&file, "file", "", "action config file path")
		configureCmd.PersistentFlags().StringVar(&prefix, "prefix", "", "action config prefix")
		configureCmd.PersistentFlags().StringArrayVar(&excludes, "excludes", []string{}, "action config prefix")
		rootCmd.AddCommand(configureCmd)
}

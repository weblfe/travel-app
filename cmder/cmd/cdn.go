package cmd

import (
	"github.com/spf13/cobra"
	"github.com/weblfe/travel-app/cmder/kernel"
)

// dockerCmd represents the docker command
var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "cdn service",
	Long:  `cdn 域名切换工具`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("cmd args:",args)
		if cdnDomain == "" {
			cmd.Println("域名为空,无法更新")
			return
		}
		var domain = kernel.NewCdnDomain()
		domain.SetDomainUrl(cdnDomain)
	},
}
var cdnDomain string

func init() {
	// Here you will define your flags and configuration settings.
	cdnCmd.PersistentFlags().StringVar(&cdnDomain, "domain", "", "set cdn domain url")
	rootCmd.AddCommand(cdnCmd)
}

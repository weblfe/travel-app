package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/weblfe/travel-app/cmder/kernel"
)

// dockerCmd represents the docker command
var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "cdn service",
	Long:  `cdn 域名切换工具`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("cmd args:", args)
		if oldCdnDomainUrl == "" {
			cmd.Println("老域名为空,无法更新")
			return
		}
		if newCdnDomainUrl == "" {
			cmd.Println("新域名为空,无法更新")
			return
		}
		var domain = kernel.NewCdnDomain()
		if domain.Replaces(oldCdnDomainUrl, newCdnDomainUrl) > 0 {
			cmd.Println(fmt.Sprintf("%s => %s", oldCdnDomainUrl, newCdnDomainUrl), "更新 cdn 域名 成功")
		}
	},
}

var (
	oldCdnDomainUrl string
	newCdnDomainUrl string
)

func init() {
	// Here you will define your flags and configuration settings.
	cdnCmd.PersistentFlags().StringVar(&oldCdnDomainUrl, "oldUrl", "", "cdn old domain url")
	cdnCmd.PersistentFlags().StringVar(&newCdnDomainUrl, "newUrl", "", "new cdn domain url")
	rootCmd.AddCommand(cdnCmd)
}

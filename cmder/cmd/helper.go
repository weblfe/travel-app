package cmd

import (
	"github.com/spf13/cobra"
)

// dockerCmd represents the docker command
var helperCmd = &cobra.Command{
	Use:   "utils",
	Short: "utils service",
	Long:  `帮助工具命令`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docker 宿主机 http api 地址
/*	helperCmd.PersistentFlags().StringArrayVar(&dockerHosts, "docker_host", []string{}, "docker host urls")
	// 手动注册
	helperCmd.PersistentFlags().StringVar(&service, "service", "", "register service for container")
	// 罗列运作中的容器
	helperCmd.PersistentFlags().BoolVar(&lists, "lists", false, "lists container service")
	// 查询 容器, 服务
	helperCmd.PersistentFlags().StringVar(&query, "query", "", "query service or container,eg: service=app&key=info")
*/
	rootCmd.AddCommand(helperCmd)
}

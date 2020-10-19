package cmd

import (
		"github.com/spf13/cobra"
		"github.com/weblfe/travel-app/cmder/kernel"
		"log"
)

var (
		dockerHosts []string
		service     string
		lists       bool
		query       string
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
		Use:   "docker",
		Short: "docker service",
		Long:  `query server docker info, docker service ,auto register docker and register docker service`,
		Run: func(cmd *cobra.Command, args []string) {
				logic := kernel.InvokerDockerService(dockerHosts, service, lists, query,endpoints,timeout)
				if err:=logic.Exec();err!=nil {
						log.Fatal(err)
				}
		},
}

func init() {

		// Here you will define your flags and configuration settings.
		// Cobra supports Persistent Flags which will work for this command
		// and all subcommands, e.g.:
		// docker 宿主机 http api 地址
		dockerCmd.PersistentFlags().StringArrayVar(&dockerHosts, "docker_host", []string{}, "docker host urls")
		// 手动注册
		dockerCmd.PersistentFlags().StringVar(&service, "service", "", "register service for container")
		// 罗列运作中的容器
		dockerCmd.PersistentFlags().BoolVar(&lists, "lists", false, "lists container service")
		// 查询 容器, 服务
		dockerCmd.PersistentFlags().StringVar(&query, "query", "", "query service or container,eg: service=app&key=info")

		rootCmd.AddCommand(dockerCmd)
}

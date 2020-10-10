package cmd

import (
		"fmt"
		"github.com/spf13/cobra"
)

var (
		dockerHosts *[]string
		service     *string
		lists       *bool
		query       *string
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
		Use:   "docker",
		Short: "docker service",
		Long:  `query server docker info, docker service ,auto register docker and register docker service`,
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(*dockerHosts)
				fmt.Println(*service)
				fmt.Println(*lists)
				fmt.Println(*query)
		},
}

func init() {

		// Here you will define your flags and configuration settings.
		// Cobra supports Persistent Flags which will work for this command
		// and all subcommands, e.g.:
		// docker 宿主机 http api 地址
		dockerHosts = dockerCmd.PersistentFlags().StringArray("docker_host", []string{}, "docker host urls")
		// 手动注册
		service = dockerCmd.PersistentFlags().String("service", "", "register service for container")
		// 罗列运作中的容器
		lists = dockerCmd.PersistentFlags().Bool("lists", false, "lists container service")

		// 查询 容器, 服务
		query = dockerCmd.PersistentFlags().String("query", "", "query service or container,eg: service=app&key=info")

		rootCmd.AddCommand(dockerCmd)
}

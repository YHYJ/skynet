/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-20 09:52:25

Description: 执行子命令 'http'
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/skynet/cli"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start an http server",
	Long:  `Start an http server and manage its life cycle.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		portFlag, _ := cmd.Flags().GetInt("port")
		dirFlag, _ := cmd.Flags().GetString("dir")
		interactiveFlag, _ := cmd.Flags().GetBool("interactive")

		// 启动 HTTP 服务 CLI 版本
		cli.StartHttp(portFlag, dirFlag, interactiveFlag)
	},
}

func init() {
	httpCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	httpCmd.Flags().StringP("dir", "d", "PWD", "Directory to serve")
	httpCmd.Flags().BoolP("interactive", "i", false, "Start interactive mode")

	httpCmd.Flags().BoolP("help", "h", false, "help for http command")
	rootCmd.AddCommand(httpCmd)
}

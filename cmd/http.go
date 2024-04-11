/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-20 09:52:25

Description: 执行子命令 'http'
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/skynet/cli"
	"github.com/yhyj/skynet/general"
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
		interfaceFlag, _ := cmd.Flags().GetBool("interface")

		// 使用 portFlag 参数
		// 如果 portFlag 参数不在[1, 65535]范围内，则使用默认值8080
		if portFlag < 1 || portFlag > 65535 {
			portFlag = 8080
			color.Printf("%s\n", general.DangerText("Port number is invalid, using default port 8080."))
		}
		// 如果 portFlag 参数小于1024，则提示需要 root 权限
		if portFlag < 1024 {
			color.Printf("%s\n", general.DangerText("You need root privileges to listen on ports below 1024."))
			os.Exit(1)
		}

		// 处理 dirFlag 默认参数
		if dirFlag == "PWD" {
			dirFlag = general.GetVariable("PWD")
		}

		// 使用 dirFlag 参数
		if !general.FileExist(dirFlag) {
			// 如果 dirFlag 参数不是一个目录，则提示目录不存在并退出程序
			color.Error.Printf("Directory '%s' does not exist.\n", dirFlag)
			os.Exit(1)
		}
		// 获取 dirFlag 参数的绝对路径
		absDir := general.GetAbsPath(dirFlag)

		// 输出 interfaceFlag 供用户选择
		netInterfacesData, _ := general.GetNetInterfacesForCli()
		var netInterfaceNumber int
		if interfaceFlag {
			// 输出网卡信息供用户选择，输出格式为：[序号] 网卡名称 网卡IP
			for i := 1; i <= len(netInterfacesData); i++ {
				// 输出网卡信息
				color.Printf("%s %s: %s\n", general.FgGreenText("[", i, "]"), general.LightText(netInterfacesData[i]["name"]), general.LightText(netInterfacesData[i]["ip"]))
			}
			// 选择网卡编号
			color.Printf("%s", general.QuestionText("Please select the interface number: "))
			// 接收用户输入并赋值给 interfaceNumber
			fmt.Scanln(&netInterfaceNumber)
			// 如果 interfaceNumber 不在[0, len(netinterfacesData))范围内，则使用默认值
			if netInterfaceNumber < 1 || netInterfaceNumber > len(netInterfacesData) {
				netInterfaceNumber = 1
				color.Danger.Printf("Invalid interface number, using default interface <%s>\n", netInterfacesData[netInterfaceNumber]["name"])
			}
			color.Println()
		} else {
			netInterfaceNumber = 1
		}
		// 获取 address 参数
		address := netInterfacesData[netInterfaceNumber]["ip"]

		// 输出服务类型供用户选择
		serviceSlice := map[int]string{1: "Download", 2: "Upload", 3: "All"}
		var serviceNumber int
		if interfaceFlag {
			// 输出支持的服务类型供用户选择，输出格式为：[序号] 服务类型
			for i := 1; i <= len(serviceSlice); i++ {
				// 输出服务类型
				color.Printf("%s %s\n", general.FgGreenText("[", i, "]"), general.LightText(serviceSlice[i]))
			}
			// 选择服务编号
			color.Printf("%s", general.QuestionText("Please select the service number: "))
			// 接收用户输入并赋值给 serviceNumber
			fmt.Scanln(&serviceNumber)
			// 如果 serviceNumber 不在[0, len(serviceSlice))范围内，则使用默认值
			if serviceNumber < 1 || serviceNumber > len(serviceSlice) {
				serviceNumber = 3
				color.Danger.Printf("Invalid service number, using default service <%s>\n", serviceSlice[serviceNumber])
			}
			color.Println()
		} else {
			serviceNumber = 3
		}

		// 启动 http server
		switch serviceSlice[serviceNumber] {
		case "Download":
			cli.HttpDownloadServer(address, fmt.Sprint(portFlag), absDir)
		case "Upload":
			cli.HttpUploadServer(address, fmt.Sprint(portFlag), absDir)
		case "All":
			cli.HttpAllServer(address, fmt.Sprint(portFlag), absDir)
		default:
			color.Error.Println("Please select service")
			return
		}
	},
}

func init() {
	httpCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	httpCmd.Flags().StringP("dir", "d", "PWD", "Directory to serve")
	httpCmd.Flags().BoolP("interface", "i", false, "Select the net interface to use (default 0.0.0.0)")

	httpCmd.Flags().BoolP("help", "h", false, "help for http command")
	rootCmd.AddCommand(httpCmd)
}

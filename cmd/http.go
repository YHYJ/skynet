/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-20 09:52:25

Description: 程序子命令'http'时执行
*/

package cmd

import (
	"fmt"
	"os"

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

		// 使用portFlag参数
		// 如果portFlag参数不在[1, 65535]范围内，则使用默认值8080
		if portFlag < 1 || portFlag > 65535 {
			portFlag = 8080
			fmt.Printf(general.ErrorBaseFormat, "Port number is invalid, using default port 8080.")
		}
		// 如果portFlag参数小于1024，则提示需要root权限
		if portFlag < 1024 {
			fmt.Printf(general.ErrorBaseFormat, "You need root privileges to listen on ports below 1024.")
		}

		// 处理dirFlag默认参数
		if dirFlag == "PWD" {
			dirFlag = general.GetVariable("PWD")
		}

		// 使用dirFlag参数
		if !general.FileExist(dirFlag) {
			// 如果dirFlag参数不是一个目录，则提示目录不存在并退出程序
			fmt.Printf(general.ErrorBaseFormat, "Directory does not exist.")
			os.Exit(1)
		}
		// 获取dirFlag参数的绝对路径
		absDir := general.GetAbsPath(dirFlag)

		// 输出interfaceFlag供用户选择
		netInterfacesData, _ := cli.GetNetInterfaces()
		var netInterfaceNumber int
		if interfaceFlag {
			// 输出网卡信息供用户选择，输出格式为：[序号] 网卡名称 网卡IP
			for i := 1; i <= len(netInterfacesData); i++ {
				// 输出网卡信息
				fmt.Printf(general.SliceTraverseSuffixFormat, fmt.Sprintf("[%d]", i), fmt.Sprintf(" %s: ", netInterfacesData[i]["name"]), netInterfacesData[i]["ip"])
			}
			// 选择网卡编号
			fmt.Printf(general.AskFormat, "Please select the net interface: ")
			// 接收用户输入并赋值给interfaceNumber
			fmt.Scanln(&netInterfaceNumber)
			// 如果interfaceNumber不在[0, len(netinterfacesData))范围内，则使用默认值
			if netInterfaceNumber < 1 || netInterfaceNumber > len(netInterfacesData) {
				netInterfaceNumber = 1
				fmt.Printf(general.ErrorBaseFormat, fmt.Sprintf("Invalid interface number, using default interface <%s>", netInterfacesData[netInterfaceNumber]["name"]))
			}
			fmt.Println()
		} else {
			netInterfaceNumber = 1
		}
		// 获取address参数
		address := netInterfacesData[netInterfaceNumber]["ip"]

		// 输出服务类型供用户选择
		serviceSlice := map[int]string{1: "Download", 2: "Upload", 3: "All"}
		var serviceNumber int
		if interfaceFlag {
			// 输出支持的服务类型供用户选择，输出格式为：[序号] 服务类型
			for i := 1; i <= len(serviceSlice); i++ {
				// 输出服务类型
				fmt.Printf(general.SliceTraverseSuffixFormat, fmt.Sprintf("[%d]", i), " ", serviceSlice[i])
			}
			// 选择服务编号
			fmt.Printf(general.AskFormat, "Please select the service type: ")
			// 接收用户输入并赋值给serviceNumber
			fmt.Scanln(&serviceNumber)
			// 如果serviceNumber不在[0, len(serviceSlice))范围内，则使用默认值
			if serviceNumber < 1 || serviceNumber > len(serviceSlice) {
				serviceNumber = 1
				fmt.Printf(general.ErrorBaseFormat, fmt.Sprintf("Invalid service number, using default service <%s>", serviceSlice[serviceNumber]))
			}
			fmt.Println()
		} else {
			serviceNumber = 1
		}

		// 启动http server
		switch serviceSlice[serviceNumber] {
		case "Download":
			cli.HttpDownloadServer(address, fmt.Sprint(portFlag), absDir)
		case "Upload":
			cli.HttpUploadServer(address, fmt.Sprint(portFlag), absDir)
		case "All":
			cli.HttpAllServer(address, fmt.Sprint(portFlag), absDir)
		default:
			fmt.Printf(general.ErrorBaseFormat, "Please select service")
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

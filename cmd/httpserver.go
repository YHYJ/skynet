/*
File: httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-04-20 09:52:25

Description: 程序子命令'httpserver'时执行
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yhyj/skynet/function"
)

// httpserverCmd represents the httpserver command
var httpserverCmd = &cobra.Command{
	Use:   "httpserver",
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
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "Port number is invalid, using default port 8080.")
		}
		// 如果portFlag参数小于1024，则提示需要root权限
		if portFlag < 1024 {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "You need root privileges to listen on ports below 1024.")
		}

		// 处理dirFlag默认参数
		if dirFlag == "PWD" {
			dirFlag = function.GetVariable("PWD")
		}

		// 使用dirFlag参数
		if !function.FileExist(dirFlag) {
			// 如果dirFlag参数不是一个目录，则提示目录不存在并退出程序
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "Directory does not exist.")
			os.Exit(1)
		}
		// 获取dirFlag参数的绝对路径
		absDir := function.GetAbsPath(dirFlag)

		// 输出interfaceFlag供用户选择
		netInterfacesData, _ := function.GetNetInterfaces()
		var netInterfaceNumber int
		if interfaceFlag {
			// 输出网卡信息供用户选择，输出格式为：[序号] 网卡名称 网卡IP
			for i := 1; i <= len(netInterfacesData); i++ {
				// 输出网卡信息
				fmt.Printf("\x1b[36;1m[%d]\x1b[0m %s: %s\n", i, netInterfacesData[i]["name"], netInterfacesData[i]["ip"])
			}
			// 选择网卡编号
			fmt.Printf("\n\x1b[34;1m%s\x1b[0m", "Please select the net interface to use: ")
			// 接收用户输入并赋值给interfaceNumber
			fmt.Scanln(&netInterfaceNumber)
			// 如果interfaceNumber不在[0, len(netinterfacesData))范围内，则使用默认值0
			if netInterfaceNumber < 1 || netInterfaceNumber > len(netInterfacesData) {
				netInterfaceNumber = 1
				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "Interface number is invalid, using default interface. ")
			}
		} else {
			netInterfaceNumber = 1
		}

		// 获取address参数
		address := netInterfacesData[netInterfaceNumber]["ip"]

		// 启动http server
		fmt.Printf("\n\x1b[32;1mStarting http server at %s ...\x1b[0m\n", absDir)
		url := fmt.Sprintf("http://%s:%v", address, portFlag)
		fmt.Printf("\x1b[32;1mServing HTTP on %s port %v (%s).\x1b[0m\n", address, portFlag, url)
		// 输出二维码
		codeString, err := function.QrCodeString(url)
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}

		// 显示服务停止快捷键
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.")
		function.HttpServer(address, fmt.Sprint(portFlag), dirFlag)
	},
}

func init() {
	httpserverCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	httpserverCmd.Flags().StringP("dir", "d", "PWD", "Directory to serve")
	httpserverCmd.Flags().BoolP("interface", "i", false, "Select the net interface to use (default 0.0.0.0)")

	httpserverCmd.Flags().BoolP("help", "h", false, "help for httpserver command")
	rootCmd.AddCommand(httpserverCmd)
}

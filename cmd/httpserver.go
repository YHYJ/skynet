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
			function.TemplateData = map[string]interface{}{"DefaultPort": portFlag}
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", function.Translate(function.Localizer, "UseDefaultPort", function.TemplateData))
		}
		// 如果portFlag参数小于1024，则提示需要root权限并退出程序
		if portFlag < 1024 {
			function.TemplateData = map[string]interface{}{"RootPort": 1024}
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", function.Translate(function.Localizer, "PortNeedPrivilege", function.TemplateData))
			return
		}

		// 使用dirFlag参数
		// 如果dirFlag参数不是一个目录，则提示目录不存在并退出程序
		if !function.FileExist(dirFlag) {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", function.Translate(function.Localizer, "FolderNotExist", function.TemplateData))
			return
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
			fmt.Printf("\n\x1b[34;1m%s\x1b[0m", function.Translate(function.Localizer, "SelectNetInterface", function.TemplateData))
			// 接收用户输入并赋值给interfaceNumber
			fmt.Scanln(&netInterfaceNumber)
			// 如果interfaceNumber不在[0, len(netinterfacesData))范围内，则使用默认值0
			if netInterfaceNumber < 1 || netInterfaceNumber > len(netInterfacesData) {
				netInterfaceNumber = 1
				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", function.Translate(function.Localizer, "NetInterfaceUnavailable", function.TemplateData))
			}
		} else {
			netInterfaceNumber = 1
		}

		// 获取address参数
		address := netInterfacesData[netInterfaceNumber]["ip"]

		// 启动http server
		function.TemplateData = map[string]interface{}{"Address": address, "Port": portFlag}
		fmt.Printf("\n\x1b[32;1m%s\x1b[0m\n", function.Translate(function.Localizer, "ServiceInformation", function.TemplateData))
		function.TemplateData = map[string]interface{}{"Folder": absDir}
		fmt.Printf("\x1b[32;1m%s\x1b[0m\n", function.Translate(function.Localizer, "StartService", function.TemplateData))
		function.TemplateData = map[string]interface{}{"ActionKey": "Ctrl+C"}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", function.Translate(function.Localizer, "HowToStopService", function.TemplateData))
		function.HttpServer(address, fmt.Sprint(portFlag), dirFlag)
	},
}

func init() {
	httpserverCmd.Flags().IntP("port", "p", 8080, function.Translate(function.Localizer, "HelpInfoPortFlag", function.TemplateData))
	httpserverCmd.Flags().StringP("dir", "d", ".", function.Translate(function.Localizer, "HelpInfoDirFlag", function.TemplateData))
	httpserverCmd.Flags().BoolP("interface", "i", false, function.Translate(function.Localizer, "HelpInfoInterfaceFlag", function.TemplateData))

	httpserverCmd.Flags().BoolP("help", "h", false, function.Translate(function.Localizer, "HelpInfoHelpFlag", function.TemplateData))
	rootCmd.AddCommand(httpserverCmd)
}

/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 14:42:12

Description: 子命令 'http' 的实现
*/

package cli

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/yhyj/skynet/general"
)

// StartHttp 启动 HTTP 服务
//
// 参数：
//   - port: 服务端口
//   - dir: 服务目录
//   - interactive: 交互模式
func StartHttp(port int, dir string, interactive bool) {
	// 如果 port 范围不在 [1, 65535] 内，则使用默认值 8080
	if port < 1 || port > 65535 {
		port = 8080
		color.Printf("%s\n", general.DangerText("Port number is invalid, using default port 8080."))
	}
	// 如果 port 小于 1024，则提示需要 root 权限
	if port < 1024 {
		color.Printf("%s\n", general.DangerText("You need root privileges to listen on ports below 1024."))
		os.Exit(1)
	}

	if dir == "PWD" {
		dir = general.GetVariable("PWD")
	}
	// 使用 dir 参数
	if !general.FileExist(dir) {
		// 如果 dir 参数不是一个目录，则提示目录不存在并退出程序
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> No such file or directory: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), dir)
		os.Exit(1)
	}
	// 获取 dir 参数的绝对路径
	absDir := general.GetAbsPath(dir)

	// 获取 CLI 适用格式的网卡信息
	netInterfacesData, _ := general.GetNetInterfacesForCLI()
	// 网卡编号
	var netInterfaceNumber int
	// 可用服务类型
	serviceSlice := map[int]string{1: "Download", 2: "Upload", 3: "All"}
	// 服务类型编号
	var serviceNumber int

	if interactive { // 交互模式
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
			color.Warn.Printf("Invalid interface number, using default interface <%s>\n", netInterfacesData[netInterfaceNumber]["name"])
		}
		color.Println()

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
			color.Warn.Printf("Invalid service number, using default service <%s>\n", serviceSlice[serviceNumber])
		}
		color.Println()
	} else { // 默认模式
		netInterfaceNumber = 1
		serviceNumber = 3
	}
	// 获取 address 参数
	address := netInterfacesData[netInterfaceNumber]["ip"]

	// 启动 http server
	switch serviceSlice[serviceNumber] {
	case "Download":
		general.HttpDownloadServerForCLI(address, color.Sprint(port), absDir)
	case "Upload":
		general.HttpUploadServerForCLI(address, color.Sprint(port), absDir)
	case "All":
		general.HttpAllServerForCLI(address, color.Sprint(port), absDir)
	default:
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> Unable to start service: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), serviceSlice[serviceNumber])
		return
	}
}

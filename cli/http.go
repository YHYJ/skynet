/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 14:42:12

Description: 子命令`http`功能实现
*/

package cli

import (
	"fmt"
	"net"
	"net/http"

	"github.com/yhyj/skynet/general"
)

// 启动HTTP服务
func HttpServer(address string, port string, dir string) {
	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		fmt.Printf("\n\x1b[32;1mStarting http server at %s\x1b[0m\n", dir) // 服务地址
		fmt.Printf("\x1b[32;1mHTTP url is %s\x1b[0m\n", url)               // URL
		codeString, err := general.QrCodeString(url)                       // 二维码
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.") // 服务停止快捷键

		// 创建请求处理器
		handler := http.FileServer(http.Dir(dir))
		// 启动服务器
		if err := http.Serve(listener, handler); err == http.ErrServerClosed {
			fmt.Printf("HTTP Server closed\n")
		} else if err != nil {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		}
	}
}

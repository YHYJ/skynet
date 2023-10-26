/*
File: httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 14:42:12

Description: 子命令`httpserver`功能函数
*/

package function

import (
	"fmt"
	"net"
	"net/http"
)

// 启动HTTP服务（Terminal使用）
func HttpServer(address string, port string, dir string) {
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	} else {
		// 成功后输出各种信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		fmt.Printf("\n\x1b[32;1mStarting http server at %s\x1b[0m\n", dir)                   // 服务地址
		fmt.Printf("\x1b[32;1mHTTP url is %s\x1b[0m\n", url) // URL
		codeString, err := QrCodeString(url)                                                 // 二维码
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.") // 服务停止快捷键

		if err := http.Serve(listener, http.FileServer(http.Dir(dir))); err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		}
	}
}

// 启动HTTP服务（GUI使用）
func HttpServerForGui(address string, port string, dir string) *http.Server {
	// 创建一个HTTP服务器
	server := &http.Server{
		Addr:    address + ":" + port,           // 指定服务器侦听的TCP地址
		Handler: http.FileServer(http.Dir(dir)), // 调用的处理程序
	}

	// 启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		}
	}()

	return server
}

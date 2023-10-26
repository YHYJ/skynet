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
	"net/http"
)

// 启动HTTP服务（Terminal使用）
func HttpServer(address string, port string, dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(address+":"+port, nil)
}

// 启动HTTP服务（GUI使用）
func HttpServerForGui(address string, port string, dir string) *http.Server {
	// 创建一个HTTP服务器
	server := &http.Server{
		Addr: address + ":" + port,
	}

	// 定义一个处理函数
	http.Handle("/", http.FileServer(http.Dir(dir)))
	// 启动HTTP服务器
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		} else if err == nil {
			fmt.Printf("HTTP server started\n")
		}
	}()

	return server
}

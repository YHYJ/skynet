/*
File: define_httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-31 16:54:05

Description: 用于GUI的HTTP服务
*/

package gui

import (
	"fmt"
	"net"
	"net/http"
)

// 启动HTTP服务
func HttpServer(address string, port string, dir string) (*http.Server, error) {
	// 创建一个HTTP服务器
	server := &http.Server{
		Addr:    address + ":" + port,           // 指定服务器侦听的TCP地址
		Handler: http.FileServer(http.Dir(dir)), // 调用的处理程序
	}

	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动HTTP服务器
		go func() {
			if err := server.Serve(listener); err == http.ErrServerClosed {
				fmt.Printf("HTTP Server closed\n")
			} else if err != nil {
				fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
			}
		}()
	}

	return server, nil
}

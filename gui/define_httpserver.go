/*
File: define_httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-31 16:54:05

Description: 用于GUI的HTTP服务
*/

package gui

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
)

// 启动HTTP服务
func HttpServer(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	// 创建一个HTTP服务器
	server := &http.Server{
		Handler: http.DefaultServeMux, // 调用的处理程序
	}

	// 在DefaultServeMux中注册给定模式的处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 列出文件夹中的所有文件，并提供下载链接
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Println("Error reading download directory:", err)
		}

		templateString := `
		<!doctype html>
		<html>
			<head><title>Download</title></head>
			<body>
				<h1>File Download</h1>
				<hr>
				<ul>
					{{range .}}
						<li><a href="/download/{{.Name}}">{{.Name}}</a></li>
					{{end}}
				</ul>
			</body>
		</html>
		`
		newTemplate, _ := template.New("download").Parse(templateString)
		newTemplate.Execute(w, files)
	})

	// 在DefaultServeMux中注册给定模式的处理程序
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动HTTP服务器
		go func() {
			if err := server.Serve(listener); err == http.ErrServerClosed {
				log.Printf("\x1b[33;1mHTTP Server closed\x1b[0m\n")
			} else if err != nil {
				log.Printf("\x1b[31;1mHTTP server error\x1b[0m: \x1b[31m%s\x1b[0m\n", "123")
			}
		}()
	}

	return server, nil
}

/*
File: http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-31 16:54:05

Description: 子命令 'gui' 的 httpserver 功能实现
*/

package gui

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yhyj/skynet/general"
)

// HttpDownloadServer 启动 HTTP 下载服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpDownloadServer(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !general.FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 加锁，控制对 HTTP 服务器和路由注册的并发访问
	// 确保只有一个 goroutine 能够启动 HTTP 服务器和注册路由，防止多次重复注册相同的路由
	general.ServerMutex.Lock()
	defer general.ServerMutex.Unlock()

	// 创建路由
	general.ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	general.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 列出文件夹中的所有文件，并提供下载链接
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(w, "Error reading download directory: %s", err)
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
	// 注册给定模式的处理程序
	general.ServeMux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

	// 创建 HTTP 服务器
	general.HttpServer = &http.Server{
		Handler: general.ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := general.HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(general.FgYellowText("HTTP Server closed"))
				general.ServeMux = nil
				general.HttpServer = nil
			} else if err != nil {
				log.Printf("%s %s\n", general.FgRedText("HTTP server error:"), general.ErrorText(err))
				general.ServeMux = nil
				general.HttpServer = nil
			}
		}()
	}

	return general.HttpServer, nil
}

// HttpUploadServer 启动 HTTP 上传服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpUploadServer(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !general.FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 创建路由
	general.ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	general.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 解析表单
			err := r.ParseMultipartForm(100 << 20) // 限制内存最多存储100MB，超出的部分保存到磁盘
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			file, handler, err := r.FormFile("file") // 获取上传文件
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			// 创建文件保存到服务目录
			targetFile, err := os.Create(filepath.Join(dir, handler.Filename))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFile.Close()

			// 将上传文件内容复制到新文件
			_, err = io.Copy(targetFile, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// JS 显示弹窗通知
			js := fmt.Sprintf(`
			<script>
				alert("File uploaded successfully\n%s");
				window.location.href = '/upload';
			</script>
			`, handler.Filename)
			fmt.Fprintln(w, js)
		} else {
			// 显示文件上传表单
			templateString := `
			<!doctype html>
			<html>
				<head><title>Upload</title></head>
				<body>
					<h1>File Upload</h1>
					<hr><br>
					<form action="/upload" method="post" enctype="multipart/form-data">
						<input type="file" name="file">
						<input type="submit" value="Upload">
					</form>
				</body>
			</html>
			`
			newTemplate, _ := template.New("upload").Parse(templateString)
			newTemplate.Execute(w, nil)
		}
	})

	// 创建 HTTP 服务器
	general.HttpServer = &http.Server{
		Handler: general.ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := general.HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(general.FgYellowText("HTTP Server closed"))
				general.ServeMux = nil
				general.HttpServer = nil
			} else if err != nil {
				log.Printf("%s %s\n", general.FgRedText("HTTP server error:"), general.ErrorText(err))
				general.ServeMux = nil
				general.HttpServer = nil
			}
		}()
	}

	return general.HttpServer, nil
}

// HttpAllServer 启动 HTTP 所有服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpAllServer(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !general.FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 创建路由
	general.ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	general.ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 根路径上显示一个链接到 /upload 页面
		templateString := `
		<!doctype html>
		<html>
			<head><title>File Service</title></head>
			<body>
				<h1>Welcome to the File Service</h1>
				<hr>
				<a href="/upload-service">File Upload</a><br>
				<a href="/download-service">File Download</a>
			</body>
		</html>
		`
		newTemplate, _ := template.New("root").Parse(templateString)
		newTemplate.Execute(w, nil)
	})
	general.ServeMux.HandleFunc("/upload-service", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 解析表单
			err := r.ParseMultipartForm(100 << 20) // 限制内存最多存储100MB，超出的部分保存到磁盘
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			file, handler, err := r.FormFile("file") // 获取上传文件
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			// 创建文件保存到服务目录
			targetFile, err := os.Create(filepath.Join(dir, handler.Filename))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer targetFile.Close()

			// 将上传文件内容复制到新文件
			_, err = io.Copy(targetFile, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// JS 显示弹窗通知
			js := fmt.Sprintf(`
			<script>
				alert("File uploaded successfully\n%s");
				window.location.href = '/upload-service';
			</script>
			`, handler.Filename)
			fmt.Fprintln(w, js)
		} else {
			// 显示文件上传表单
			templateString := `
			<!doctype html>
			<html>
				<head><title>Upload</title></head>
				<body>
					<h1>File Upload</h1>
					<hr>
					<a href="/">Back to Home Page</a>
					<a href="/download-service">Go to Download Page</a>
					<br><br>
					<form action="/upload-service" method="post" enctype="multipart/form-data">
						<input type="file" name="file">
						<input type="submit" value="Upload">
					</form>
				</body>
			</html>
			`
			newTemplate, _ := template.New("upload").Parse(templateString)
			newTemplate.Execute(w, nil)
		}
	})
	general.ServeMux.HandleFunc("/download-service", func(w http.ResponseWriter, r *http.Request) {
		// 列出文件夹中的所有文件，并提供下载链接
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(w, "Error reading download directory: %s", err)
		}

		templateString := `
		<!doctype html>
		<html>
			<head><title>Download</title></head>
			<body>
				<h1>File Download</h1>
				<hr>
				<a href="/">Back to Home Page</a>
				<a href="/upload-service">Go to Upload Page</a>
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
	// 注册给定模式的处理程序
	general.ServeMux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

	// 创建 HTTP 服务器
	general.HttpServer = &http.Server{
		Handler: general.ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := general.HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(general.FgYellowText("HTTP Server closed"))
				general.ServeMux = nil
				general.HttpServer = nil
			} else if err != nil {
				log.Printf("%s %s\n", general.FgRedText("HTTP server error:"), general.ErrorText(err))
				general.ServeMux = nil
				general.HttpServer = nil
			}
		}()
	}

	return general.HttpServer, nil
}

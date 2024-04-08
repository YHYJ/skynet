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
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/yhyj/skynet/general"
)

// HttpDownloadServer 启动 HTTP 下载服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpDownloadServer(address string, port string, dir string) {
	method := "Download"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		color.Error.Println(err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", general.SuccessText(method), general.FgCyan(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", general.FgBlue(url))                                          // URL
		codeString, err := general.QrCodeString(url)                                                           // 二维码
		if err != nil {
			color.Error.Println(err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", general.CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

		// 创建请求处理器
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			newTemplate, _ := template.New(strings.ToLower(method)).Parse(templateString)
			newTemplate.Execute(w, files)
		})
		http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

		// 启动服务器
		if err := http.Serve(listener, nil); err == http.ErrServerClosed {
			color.Printf("HTTP Server closed\n")
		} else if err != nil {
			color.Error.Printf("%s: %s\n", "HTTP server error", err)
		}
	}
}

// HttpUploadServer 启动 HTTP 上传服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpUploadServer(address string, port string, dir string) {
	method := "Upload"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		color.Error.Println(err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", general.SuccessText(method), general.FgCyan(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", general.FgBlue(url))                                          // URL
		codeString, err := general.QrCodeString(url)                                                           // 二维码
		if err != nil {
			color.Error.Println(err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", general.CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

		// 在 DefaultServeMux 中注册给定模式的处理函数
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// 解析表单
				err := r.ParseMultipartForm(10 << 20) // 限制上传文件大小
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				file, handler, err := r.FormFile("file")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer file.Close()

				// 创建文件保存到 uploads 文件夹
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
				// 返回包含 JavaScript 的响应以显示弹窗通知
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
				newTemplate, _ := template.New(strings.ToLower(method)).Parse(templateString)
				newTemplate.Execute(w, nil)
			}
		})

		// 启动服务器
		if err := http.Serve(listener, nil); err == http.ErrServerClosed {
			color.Printf("HTTP Server closed\n")
		} else if err != nil {
			color.Error.Printf("%s: %s\n", "HTTP server error", err)
		}
	}
}

// HttpAllServer 启动所有 HTTP 服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpAllServer(address string, port string, dir string) {
	method := "All"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		color.Error.Println(err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", general.SuccessText(method), general.FgCyan(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", general.FgBlue(url))                                          // URL
		codeString, err := general.QrCodeString(url)                                                           // 二维码
		if err != nil {
			color.Error.Println(err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", general.CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

		// 在 DefaultServeMux 中注册给定模式的处理函数
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 根路径上显示一个链接到 /upload 页面
			templateString := `
			<!doctype html>
			<html>
				<head><title>File Service</title></head>
				<body>
					<h1>Welcome to the File Service</h1>
					<a href="/upload">File Upload</a><br>
					<a href="/download">File Download</a>
				</body>
			</html>
			`
			newTemplate, _ := template.New("root").Parse(templateString)
			newTemplate.Execute(w, nil)
		})
		// 启动 Upload 服务器
		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// 解析表单
				err := r.ParseMultipartForm(10 << 20) // 限制上传文件大小
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				file, handler, err := r.FormFile("file")
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				defer file.Close()

				// 创建文件保存到 uploads 文件夹
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
				// 返回包含 JavaScript 的响应以显示弹窗通知
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
						<a href="/">Back to Main Page</a>
						<a href="/download">Go to Download Page</a>
						<br><br>
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
		// 启动 Download 服务器
		http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
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
					<a href="/">Back to Main Page</a>
					<a href="/upload">Go to Upload Page</a>
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

		// 启动服务器
		if err := http.Serve(listener, nil); err == http.ErrServerClosed {
			color.Printf("HTTP Server closed\n")
		} else if err != nil {
			color.Error.Printf("%s: %s\n", "HTTP server error", err)
		}
	}
}

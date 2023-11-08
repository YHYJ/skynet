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
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yhyj/skynet/general"
)

// HttpDownloadServer 启动 HTTP 下载服务
func HttpDownloadServer(address string, port string, dir string) {
	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		fmt.Printf("\x1b[32;1mStarting http server [Download] at '%s'\x1b[0m\n", dir) // 服务地址
		fmt.Printf("\x1b[32;1mHTTP url is %s\x1b[0m\n", url)                          // URL
		codeString, err := general.QrCodeString(url)                                  // 二维码
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.") // 服务停止快捷键

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
			newTemplate, _ := template.New("download").Parse(templateString)
			newTemplate.Execute(w, files)
		})
		http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

		// 启动服务器
		if err := http.Serve(listener, nil); err == http.ErrServerClosed {
			fmt.Printf("HTTP Server closed\n")
		} else if err != nil {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		}
	}
}

// HttpUploadServer 启动 HTTP 上传服务
func HttpUploadServer(address string, port string, dir string) {
	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		fmt.Printf("\x1b[32;1mStarting http server [Upload] at '%s'\x1b[0m\n", dir) // 服务地址
		fmt.Printf("\x1b[32;1mHTTP url is %s\x1b[0m\n", url)                        // URL
		codeString, err := general.QrCodeString(url)                                // 二维码
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.") // 服务停止快捷键

		// 在DefaultServeMux中注册给定模式的处理函数
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

				// 创建文件保存到uploads文件夹
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
				// 返回包含JavaScript的响应以显示弹窗通知
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

		// 启动服务器
		if err := http.Serve(listener, nil); err == http.ErrServerClosed {
			fmt.Printf("HTTP Server closed\n")
		} else if err != nil {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		}
	}
}

// HttpAllServer 启动所有 HTTP 服务
func HttpAllServer(address string, port string, dir string) {
	// 创建TCP监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	} else {
		// 成功后输出服务信息
		url := fmt.Sprintf("http://%s:%v", address, port)
		fmt.Printf("\x1b[32;1mStarting http server [All] at '%s'\x1b[0m\n", dir) // 服务地址
		fmt.Printf("\x1b[32;1mHTTP url is %s\x1b[0m\n", url)                     // URL
		codeString, err := general.QrCodeString(url)                             // 二维码
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		} else {
			fmt.Printf("\n%s", codeString)
		}
		fmt.Printf("\n\x1b[33;1m%s\x1b[0m\n", "Press Ctrl+C to stop.") // 服务停止快捷键

		// 在DefaultServeMux中注册给定模式的处理函数
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 根路径上显示一个链接到/upload页面
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

				// 创建文件保存到uploads文件夹
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
				// 返回包含JavaScript的响应以显示弹窗通知
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
			fmt.Printf("HTTP Server closed\n")
		} else if err != nil {
			fmt.Printf("HTTP server error: \x1b[31;1m%s\x1b[0m\n", err)
		}
	}
}

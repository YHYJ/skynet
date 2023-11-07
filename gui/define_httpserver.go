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
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// 启动HTTP下载服务
func HttpDownloadServer(address string, port string, dir string) (*http.Server, error) {
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
			log.Println("Error reading download directory: ", err)
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
				log.Printf("\x1b[31;1mHTTP server error\x1b[0m: \x1b[31m%s\x1b[0m\n", err)
			}
		}()
	}

	return server, nil
}

// 启动HTTP上传服务
func HttpUploadServer(address string, port string, dir string) (*http.Server, error) {
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
				log.Printf("\x1b[31;1mHTTP server error\x1b[0m: \x1b[31m%s\x1b[0m\n", err)
			}
		}()
	}

	return server, nil
}

// 启动HTTP下载/上传服务
func HttpDownloadUploadServer(address string, port string, dir string) (*http.Server, error) {
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
			log.Println("Error reading uploads directory: ", err)
			return
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
				log.Printf("\x1b[31;1mHTTP server error\x1b[0m: \x1b[31m%s\x1b[0m\n", err)
			}
		}()
	}

	return server, nil
}

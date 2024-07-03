/*
File: define_http.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 16:24:41

Description: 管理 HTTP 服务
*/

package general

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/gookit/color"
)

var (
	HttpServer  *http.Server   // HTTP 服务
	ServerMutex sync.Mutex     // 互斥锁
	ServeMux    *http.ServeMux // 路由
)

// HttpDownloadServerForCLI 启动 HTTP 下载服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpDownloadServerForCLI(address string, port string, dir string) {
	method := "Download"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
	} else {
		// 成功后输出服务信息
		url := color.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", SuccessText(method), FgCyanText(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", FgBlueText(url))                                  // URL
		codeString, err := QrCodeString(url)                                                       // 二维码
		if err != nil {
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

		// 创建请求处理器
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// 列出文件夹中的所有文件，并提供下载链接
			files, err := os.ReadDir(dir)
			if err != nil {
				color.Fprintf(w, "Error reading download directory: %s", err)
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
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		}
	}
}

// HttpUploadServerForCLI 启动 HTTP 上传服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpUploadServerForCLI(address string, port string, dir string) {
	method := "Upload"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
	} else {
		// 成功后输出服务信息
		url := color.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", SuccessText(method), FgCyanText(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", FgBlueText(url))                                  // URL
		codeString, err := QrCodeString(url)                                                       // 二维码
		if err != nil {
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

		// 在 DefaultServeMux 中注册给定模式的处理函数
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				// 解析表单
				if err := r.ParseMultipartForm(10 << 20); err != nil { // 限制上传文件大小
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
				if _, err = io.Copy(targetFile, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// 返回包含 JavaScript 的响应以显示弹窗通知
				js := color.Sprintf(`
				<script>
					alert("File uploaded successfully\n%s");
					window.location.href = '/upload';
				</script>
				`, handler.Filename)
				color.Fprintln(w, js)
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
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		}
	}
}

// HttpAllServerForCLI 启动所有 HTTP 服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
func HttpAllServerForCLI(address string, port string, dir string) {
	method := "All"
	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
	} else {
		// 成功后输出服务信息
		url := color.Sprintf("http://%s:%v", address, port)
		color.Info.Tips("Starting HTTP [%s] server at '%s'", SuccessText(method), FgCyanText(dir)) // 服务地址
		color.Info.Tips("HTTP server url is %s", FgBlueText(url))                                  // URL
		codeString, err := QrCodeString(url)                                                       // 二维码
		if err != nil {
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		} else {
			color.Printf("\n%s\n", codeString)
		}
		color.Printf("%s\n", CommentText("Press Ctrl+C to stop.")) // 服务停止快捷键

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
				if err := r.ParseMultipartForm(10 << 20); err != nil { // 限制上传文件大小
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
				if _, err = io.Copy(targetFile, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// 返回包含 JavaScript 的响应以显示弹窗通知
				js := color.Sprintf(`
				<script>
					alert("File uploaded successfully\n%s");
					window.location.href = '/upload';
				</script>
				`, handler.Filename)
				color.Fprintln(w, js)
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
				color.Fprintf(w, "Error reading download directory: %s", err)
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
			fileName, lineNo := GetCallerInfo()
			color.Printf("%s %s %s\n", DangerText(ErrorInfoFlag), SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		}
	}
}

// HttpDownloadServerForGUI 启动 HTTP 下载服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpDownloadServerForGUI(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 加锁，控制对 HTTP 服务器和路由注册的并发访问
	// 确保只有一个 goroutine 能够启动 HTTP 服务器和注册路由，防止多次重复注册相同的路由
	ServerMutex.Lock()
	defer ServerMutex.Unlock()

	// 创建路由
	ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 列出文件夹中的所有文件，并提供下载链接
		files, err := os.ReadDir(dir)
		if err != nil {
			color.Fprintf(w, "Error reading download directory: %s", err)
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
	ServeMux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

	// 创建 HTTP 服务器
	HttpServer = &http.Server{
		Handler: ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(FgYellowText("HTTP Server closed"))
				ServeMux = nil
				HttpServer = nil
			} else if err != nil {
				log.Printf("%s\n", DangerText("HTTP server error: ", err))
				ServeMux = nil
				HttpServer = nil
			}
		}()
	}

	return HttpServer, nil
}

// HttpUploadServerForGUI 启动 HTTP 上传服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpUploadServerForGUI(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 创建路由
	ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 解析表单
			if err := r.ParseMultipartForm(100 << 20); err != nil { // 限制内存最多存储100MB，超出的部分保存到磁盘
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
			if _, err = io.Copy(targetFile, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// JS 显示弹窗通知
			js := color.Sprintf(`
			<script>
				alert("File uploaded successfully\n%s");
				window.location.href = '/upload';
			</script>
			`, handler.Filename)
			color.Fprintln(w, js)
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
	HttpServer = &http.Server{
		Handler: ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(FgYellowText("HTTP Server closed"))
				ServeMux = nil
				HttpServer = nil
			} else if err != nil {
				log.Printf("%s\n", DangerText("HTTP server error: ", err))
				ServeMux = nil
				HttpServer = nil
			}
		}()
	}

	return HttpServer, nil
}

// HttpAllServerForGUI 启动 HTTP 所有服务
//
// 参数：
//   - address: 服务地址
//   - port: 服务端口
//   - dir: 服务目录
//
// 返回：
//   - HTTP 服务器对象
//   - 错误信息
func HttpAllServerForGUI(address string, port string, dir string) (*http.Server, error) {
	// 服务启动目录不存在则创建
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 创建路由
	ServeMux = http.NewServeMux()
	// 注册给定模式的处理函数
	ServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	ServeMux.HandleFunc("/upload-service", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 解析表单
			if err := r.ParseMultipartForm(100 << 20); err != nil { // 限制内存最多存储100MB，超出的部分保存到磁盘
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
			if _, err = io.Copy(targetFile, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// JS 显示弹窗通知
			js := color.Sprintf(`
			<script>
				alert("File uploaded successfully\n%s");
				window.location.href = '/upload-service';
			</script>
			`, handler.Filename)
			color.Fprintln(w, js)
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
	ServeMux.HandleFunc("/download-service", func(w http.ResponseWriter, r *http.Request) {
		// 列出文件夹中的所有文件，并提供下载链接
		files, err := os.ReadDir(dir)
		if err != nil {
			color.Fprintf(w, "Error reading download directory: %s", err)
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
	ServeMux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(dir))))

	// 创建 HTTP 服务器
	HttpServer = &http.Server{
		Handler: ServeMux, // 调用的处理程序
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", address+":"+port)
	if err != nil {
		return nil, err
	} else {
		// 启动 HTTP 服务器
		go func() {
			if err := HttpServer.Serve(listener); err == http.ErrServerClosed {
				log.Println(FgYellowText("HTTP Server closed"))
				ServeMux = nil
				HttpServer = nil
			} else if err != nil {
				log.Printf("%s\n", DangerText("HTTP server error: ", err))
				ServeMux = nil
				HttpServer = nil
			}
		}()
	}

	return HttpServer, nil
}

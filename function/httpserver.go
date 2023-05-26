/*
File: httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 14:42:12

Description: 子命令`httpserver`功能函数
*/

package function

import "net/http"

func HttpServer(address string, port string, dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(address+":"+port, nil)
}

/*
File: httpserver.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 14:42:12

Description: 子命令`httpserver`功能函数
*/

package function

import (
	"net/http"

	"github.com/skip2/go-qrcode"
)

// 启动HTTP服务
func HttpServer(address string, port string, dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(address+":"+port, nil)
}

// 生成二维码
func QrCode(url string) (string, error) {
	code, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return "", err
	}
	codeString := code.ToSmallString(false)

	return codeString, nil
}

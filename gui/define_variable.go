/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-09 10:20:07

Description: 定义变量
*/

package gui

import (
	"net/http"
	"sync"
)

var (
	serverMutex sync.Mutex     //互斥锁
	serveMux    *http.ServeMux // 路由
	httpServer  *http.Server   // HTTP 服务
)

/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-09 10:20:07

Description: 定义变量
*/

package gui

import (
	"fmt"
	"net/http"
	"sync"
)

var (
	serverMutex sync.Mutex                                 // 互斥锁
	serveMux    *http.ServeMux                             // 路由
	httpServer  *http.Server                               // HTTP 服务
	otherNic    string                                     // 其他网络接口
	defaultNic  = fmt.Sprintf("%s - %s", "any", "0.0.0.0") // 默认网络接口
)

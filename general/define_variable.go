/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:01:45

Description: 操作变量
*/

package general

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"sync"

	"github.com/gookit/color"
)

// ---------- 代码变量

var (
	FgBlack   = color.FgBlack.Render   // 前景色 - 黑色
	FgWhite   = color.FgWhite.Render   // 前景色 - 白色
	FgGray    = color.FgGray.Render    // 前景色 - 灰色
	FgRed     = color.FgRed.Render     // 前景色 - 红色
	FgGreen   = color.FgGreen.Render   // 前景色 - 绿色
	FgYellow  = color.FgYellow.Render  // 前景色 - 黄色
	FgBlue    = color.FgBlue.Render    // 前景色 - 蓝色
	FgMagenta = color.FgMagenta.Render // 前景色 - 品红
	FgCyan    = color.FgCyan.Render    // 前景色 - 青色

	BgBlack   = color.BgBlack.Render   // 背景色 - 黑色
	BgWhite   = color.BgWhite.Render   // 背景色 - 白色
	BgGray    = color.BgGray.Render    // 背景色 - 灰色
	BgRed     = color.BgRed.Render     // 背景色 - 红色
	BgGreen   = color.BgGreen.Render   // 背景色 - 绿色
	BgYellow  = color.BgYellow.Render  // 背景色 - 黄色
	BgBlue    = color.BgBlue.Render    // 背景色 - 蓝色
	BgMagenta = color.BgMagenta.Render // 背景色 - 品红
	BgCyan    = color.BgCyan.Render    // 背景色 - 青色

	InfoText      = color.Info.Render      // Info 文本
	NoteText      = color.Note.Render      // Note 文本
	LightText     = color.Light.Render     // Light 文本
	ErrorText     = color.Error.Render     // Error 文本
	DangerText    = color.Danger.Render    // Danger 文本
	NoticeText    = color.Notice.Render    // Notice 文本
	SuccessText   = color.Success.Render   // Success 文本
	CommentText   = color.Comment.Render   // Comment 文本
	PrimaryText   = color.Primary.Render   // Primary 文本
	WarnText      = color.Warn.Render      // Warn 文本
	QuestionText  = color.Question.Render  // Question 文本
	SecondaryText = color.Secondary.Render // Secondary 文本
)

var (
	ServerMutex sync.Mutex                                 // 互斥锁
	ServeMux    *http.ServeMux                             // 路由
	HttpServer  *http.Server                               // HTTP 服务
	OtherNic    string                                     // 其他网络接口
	DefaultNic  = fmt.Sprintf("%s - %s", "any", "0.0.0.0") // 默认网络接口
)

// ---------- 环境变量

var Platform = runtime.GOOS                   // 操作系统
var Arch = runtime.GOARCH                     // 系统架构
var UserInfo, _ = GetUserInfoByName(UserName) // 用户信息
// 用户名，当程序提权运行时，使用 SUDO_USER 变量获取提权前的用户名
var UserName = func() string {
	if GetVariable("SUDO_USER") != "" {
		return GetVariable("SUDO_USER")
	}
	return GetVariable("USER")
}()

// 用来处理不同系统之间的变量名差异
var platformChart = map[string]map[string]string{
	"windows": {
		"HOME":     "USERPROFILE",  // 用户主目录路径
		"USER":     "USERNAME",     // 当前登录用户名
		"SHELL":    "ComSpec",      // 默认shell或命令提示符路径
		"PWD":      "CD",           // 当前工作目录路径
		"HOSTNAME": "COMPUTERNAME", // 计算机主机名
	},
}

// GetVariable 获取环境变量
//
// 参数：
//   - key: 变量名
//
// 返回：
//   - 变量值
func GetVariable(key string) string {
	if innerMap, exists := platformChart[Platform]; exists {
		if _, variableExists := innerMap[key]; variableExists {
			key = platformChart[Platform][key]
		}
	}
	variable := os.Getenv(key)

	return variable
}

// GetHostname 获取系统 HOSTNAME
//
// 返回：
//   - HOSTNAME 或空字符串
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SetVariable 设置环境变量
//
// 参数：
//   - key: 变量名
//   - value: 变量值
//
// 返回：
//   - 错误信息
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// GetUserInfoByName 根据用户名获取用户信息
//
// 参数：
//   - userName: 用户名
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetUserInfoByName(userName string) (*user.User, error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoById 根据 ID 获取用户信息
//
// 参数：
//   - userId: 用户 ID
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetUserInfoById(userId int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetCurrentUserInfo 获取当前用户信息
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}

/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 10:05:05

Description: 子命令`version`功能实现
*/

package general

import "fmt"

// 程序信息
const (
	Name    string = "Skynet"
	Version string = "v0.8.7"
	Project string = "github.com/yhyj/skynet"
)

// 编译信息
var (
	GitCommitHash string = "unknown"
	BuildTime     string = "unknown"
	BuildBy       string = "unknown"
)

// ProgramInfo 返回程序信息
func ProgramInfo(only bool) string {
	programInfo := fmt.Sprintf("%s\n", Version)
	if !only {
		programInfo = fmt.Sprintf("%s version: %s\nGit commit hash: %s\nBuilt on: %s\nBuilt by: %s\n", Name, Version, GitCommitHash, BuildTime, BuildBy)
	}
	return programInfo
}

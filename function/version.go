/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 10:05:05

Description: 子命令`version`功能函数
*/

package function

import "fmt"

// 程序信息
const (
	name    string = "Skynet"
	version string = "v0.6.8"
	project string = "github.com/yhyj/skynet"
)

// 编译信息
var (
	gitCommitHash string = "unknown"
	buildTime     string = "unknown"
	buildBy       string = "unknown"
)

func ProgramInfo(only bool) string {
	programInfo := fmt.Sprintf("%s\n", version)
	if !only {
		programInfo = fmt.Sprintf("%s version: %s\nGit commit hash: %s\nBuilt on: %s\nBuilt by: %s\n", name, version, gitCommitHash, buildTime, buildBy)
	}
	return programInfo
}

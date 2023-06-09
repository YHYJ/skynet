/*
File: file_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-26 11:15:42

Description: 文件操作
*/

package function

import (
	"os"
	"path/filepath"
)

// 判断文件是否存在
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 获取指定文件绝对路径
func GetAbsPath(filePath string) string {
	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	} else {
		return absPath
	}
}

/*
File: file_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-26 11:15:42

Description: 文件操作
*/

package function

import "os"

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

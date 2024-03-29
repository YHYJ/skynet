/*
File: define_gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:13:30

Description: GUI操作
*/

package gui

import (
	"os"
	"strings"

	"github.com/flopp/go-findfont"
)

// SetFont 设置 Fyne 使用的字体
//
// 返回：
//   - 错误信息
func SetFont() error {
	// 默认使用字体
	fontNames := []string{
		"pingfang",                // macOS
		"sourcehansanscn-medium",  // Linux
		"simhei", "yahei", "msyh", // Windows
	}

	// 系统可用字体
	fontPaths := findfont.List()

	for _, name := range fontNames {
		for _, path := range fontPaths {
			pathLower := strings.ToLower(path)
			// 暂时无法解析ttc文件
			if strings.Contains(pathLower, name) && !strings.HasSuffix(pathLower, ".ttc") {
				if err := os.Setenv("FYNE_FONT", path); err != nil {
					return err
				}
				return nil // 设置成功即退出
			}
		}
	}
	return nil
}

/*
File: gui_linux.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:42:59

Description: 子命令 'gui' 的实现
*/

package gui

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

// makeCustomDialog 生成自定义对话框
//
// 参数：
//   - title: 对话框标题
//   - dismiss: 按钮文本
//   - text: 对话框内容
//   - size: 对话框大小
//   - parent: 父窗口
//
// 返回：
//   - 对话框对象
func makeCustomDialog(title, dismiss, text string, size fyne.Size, parent fyne.Window) *dialog.CustomDialog {
	dialogContent := widget.NewLabel(text)     // 设置对话框内容
	dialogContent.Wrapping = fyne.TextWrapWord // 设置换行方式
	customDialog := dialog.NewCustom(title, dismiss, dialogContent, parent)
	customDialog.Resize(size)
	return customDialog
}

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

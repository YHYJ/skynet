/*
File: gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:42:59

Description: 子命令`gui`功能函数
*/

package function

import (
	"fmt"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// 生成错误提示框
func makeErrorDialog(title, dismiss, text string, size fyne.Size, parent fyne.Window) *dialog.CustomDialog {
	content := widget.NewLabel(text)     // 设置内容
	content.Wrapping = fyne.TextWrapWord // 设置换行
	errorDialog := dialog.NewCustom(title, dismiss, content, parent)
	errorDialog.Resize(size)
	return errorDialog
}

// 启动GUI
func StartGraphicalUserInterface() {
	// 创建一个新应用
	app := app.NewWithID("Skynet")
	app.SetIcon(fyne.NewStaticResource("icon", resourceIconPng.StaticContent))

	// 创建主窗口
	mainWindow := app.NewWindow("Skynet")
	mainWindow.SetMaster()                                                    // 设置为主窗口
	baseWeight, baseHeight := float32(300), mainWindow.Canvas().Size().Height // 窗口基础尺寸
	mainWindow.Resize(fyne.NewSize(baseWeight, baseHeight))                   // 设置窗口大小
	mainWindow.SetFixedSize(true)                                             // 固定窗口大小

	// 设置错误提示框尺寸
	errorDialogSize := fyne.NewSize(baseWeight-float32(20), baseHeight-float32(20))

	// 获取网卡信息
	interfaceLabel := widget.NewLabel("选择接口:")
	nics, err := GetNetInterfacesForGui()
	if err != nil {
		errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
		errorDialog.Show()
	}
	// 创建一个单选按钮组
	interfaceRadio := widget.NewRadioGroup(nics, func(selected string) {})

	// 端口选择
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("输入端口号")

	// 目录选择
	selectedFolderEntry := widget.NewEntry()
	selectedFolderEntry.SetPlaceHolder("设置服务目录")
	folderButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		// 固定文件选择对话框大小不可修改
		const width float32 = 770
		const height float32 = 481
		// 创建一个新窗口用于文件夹选择
		folderWindow := app.NewWindow("Directory Selection")
		folderWindow.Resize(fyne.NewSize(width, height))
		folderWindow.SetFixedSize(true) // 固定窗口大小
		folderWindow.CenterOnScreen()   // 居中显示

		// 显示新窗口
		folderWindow.Show()

		// 弹出文件夹选择对话框
		fileDialog := dialog.NewFolderOpen(func(dri fyne.ListableURI, err error) {
			if err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			} else if dri == nil {
				customErrText := "Startup folder not set"
				errorDialog := makeErrorDialog("Error", "Close", customErrText, errorDialogSize, mainWindow)
				errorDialog.Show()
			} else {
				// 在标签中显示选择的文件夹路径
				selectedFolderEntry.SetText(strings.Split(dri.String(), "//")[1])
			}
			// 关闭新窗口
			folderWindow.Close()
		}, folderWindow)
		fileDialog.Show()
		fileDialog.Resize(fyne.NewSize(width, height))
	})
	// 将selectedFolderEntry和folderButton放置在同一行
	dirRow := container.NewBorder(nil, nil, folderButton, nil, selectedFolderEntry)

	// 分隔线
	separator := widget.NewSeparator()

	// 服务状态显示
	statusAnimation := widget.NewProgressBarInfinite()
	statusAnimation.Stop()

	var (
		server   *http.Server   // HTTP服务
		button   *widget.Button // 服务的启动/停止按钮
		qrWindow fyne.Window    // 二维码窗口
	)

	// 参数
	var (
		// HTTP服务默认启动参数
		defaultIP   = "0.0.0.0"
		defaultPort = "8080"
		defaultDir  = GetVariable("HOME")
	)

	// 按钮状态标识
	serviceStatus := 0 // 0是服务未启动，1是服务已启动

	button = widget.NewButton("Start", func() {
		// 获取参数信息，如果参数为空则使用默认值
		selectedInterfaceIP := func() string {
			parts := strings.Split(interfaceRadio.Selected, " ")
			if len(parts) > 1 {
				return parts[len(parts)-1]
			}
			return defaultIP
		}()
		selectedPort := func() string {
			if portEntry.Text != "" {
				return portEntry.Text
			}
			return defaultPort
		}()
		selectedDir := func() string {
			if selectedFolderEntry.Text != "" {
				return selectedFolderEntry.Text
			}
			return defaultDir
		}()
		serviceUrl := fmt.Sprintf("http://%s:%v", selectedInterfaceIP, selectedPort)

		// 生成二维码
		qrCodeImage, err := QrCodeImage(serviceUrl)
		if err != nil {
			errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
			errorDialog.Show()
		}
		// 将二维码图像转换为 Fyne 图像
		qrImage := canvas.NewImageFromImage(qrCodeImage)
		// 设置图像填充模式为 ImageFillOriginal，以确保不拉伸
		qrImage.FillMode = canvas.ImageFillOriginal

		if serviceStatus == 0 {
			// 启动HTTP服务
			server, err = HttpServerForGui(selectedInterfaceIP, selectedPort, selectedDir)
			if err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			} else {
				serviceStatus = 1       // 服务已启动
				statusAnimation.Start() // 服务状态动画
				button.SetText("Stop")  // 修改按钮文字
				qrWindow = app.NewWindow("QR Code")
				// 将二维码图像添加到窗口
				// NOTE: 不能使用container.NewCenter()函数将其添加到窗口中心，否则会产生内边距
				qrWindow.SetContent(qrImage)
				// 设置窗口内边距为零以确保图像与窗口边框贴合
				qrWindow.SetPadded(false)
				qrWindow.Show() // 显示二维码窗口
				fmt.Printf("\x1b[32;1mServing HTTP on %s port %v (%s)\x1b[0m\n", selectedInterfaceIP, selectedPort, serviceUrl)
			}
		} else if serviceStatus == 1 {
			// 停止HTTP服务
			if err := server.Shutdown(nil); err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			}
			serviceStatus = 0       // 服务已停止
			statusAnimation.Stop()  // 服务状态动画
			button.SetText("Start") // 修改按钮文字
			qrWindow.Hide()         // 隐藏二维码窗口
		} else {
			customErrText := "Unknown error"
			errorDialog := makeErrorDialog("Error", "Close", customErrText, errorDialogSize, mainWindow)
			errorDialog.Show()
		}
	})

	content := container.NewVBox(
		interfaceLabel, interfaceRadio, // 网卡选择
		portEntry,       // 端口配置
		dirRow,          // 文件夹选择
		statusAnimation, // 状态显示
		separator,       // 分隔线
		button,          // 启动按钮
	)
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}

/*
File: gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 09:42:59

Description: 子命令`gui`功能实现
*/

package gui

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yhyj/skynet/general"
)

// 启动GUI
func StartGraphicalUserInterface() {
	// HTTP服务默认配置
	var (
		defaultIP   = "0.0.0.0"                                           // HTTP服务默认绑定的IP
		defaultPort = "8080"                                              // HTTP服务默认监听的端口
		defaultDir  = general.GetVariable("HOME")                         // HTTP服务默认启动路径
		serviceUrl  = fmt.Sprintf("http://%s:%s", defaultIP, defaultPort) // HTTP服务默认URL
	)

	// 界面默认配置
	var (
		portText           = "Port [1~65535]"                                 // 端口框默认文本
		selectedFolderText = fmt.Sprintf("Directory, default %s", defaultDir) // 服务启动路径框默认文本
	)

	// 预定义服务接口和小部件
	var (
		httpServer    *http.Server   // HTTP服务
		controlButton *widget.Button // 服务的启动/停止按钮
		folderButton  *widget.Button // 目录选择按钮
		urlButton     *widget.Button // 打开URL按钮
		qrButton      *widget.Button // 二维码显示/隐藏按钮
		qrWindow      fyne.Window    // 二维码窗口
	)

	// 创建一个新应用
	appInstance := app.NewWithID(general.Name)
	appInstance.SetIcon(fyne.NewStaticResource("icon", resourceIconPng.StaticContent))

	// 创建主窗口
	mainWindow := appInstance.NewWindow(fmt.Sprintf("%s - %s", general.Name, general.Version))
	mainWindow.SetMaster()                                                                           // 该窗口设为主窗口
	mainWindow.SetFixedSize(false)                                                                   // 是否固定窗口大小
	baseWeight, baseHeight := float32(len(selectedFolderText)*10), mainWindow.Canvas().Size().Height // 窗口基础尺寸
	mainWindow.Resize(fyne.NewSize(baseWeight, baseHeight))                                          // 设置窗口大小

	// 创建错误提示框尺寸
	errorDialogSize := fyne.NewSize(baseWeight-float32(20), baseHeight-float32(20))

	// 获取网卡信息
	interfaceLabel := widget.NewLabel("Select Interface:")
	nicInfos, err := GetNetInterfaces()
	if err != nil {
		errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
		errorDialog.Show()
	}
	// 创建接口选择器（单选按钮组）
	interfaceRadio := widget.NewRadioGroup(nicInfos, func(selected string) {})

	// 创建端口选择器
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder(portText)

	// 创建目录选择器标签
	selectedFolderEntry := widget.NewEntry()
	selectedFolderEntry.SetPlaceHolder(selectedFolderText)
	// 创建目录选择器
	folderButton = widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		// 固定文件选择对话框大小不可修改
		const folderWidth, folderHeight float32 = 770, 481
		// 创建一个新窗口用于文件夹选择
		folderWindow := appInstance.NewWindow("Directory Selection")
		folderWindow.Resize(fyne.NewSize(folderWidth, folderHeight))
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
		fileDialog.Resize(fyne.NewSize(folderWidth, folderHeight))
	})
	// 创建URL打开按钮
	urlButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		serviceUrlParsed, err := url.Parse(serviceUrl)
		if err != nil {
			errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
			errorDialog.Show()
		}
		appInstance.OpenURL(serviceUrlParsed)
	})
	urlButton.Disable() // 禁用URL按钮

	// 创建分隔线
	separator := widget.NewSeparator()

	// 创建服务状态显示动画
	statusAnimation := widget.NewProgressBarInfinite()
	statusAnimation.Stop()

	// 创建二维码窗口
	appDriver := appInstance.Driver()
	if drv, ok := appDriver.(desktop.Driver); ok {
		qrWindow = drv.CreateSplashWindow() // 无边框窗口
	} else {

		qrWindow = appInstance.NewWindow("QR Code") // 普通窗口
	}

	// 状态标识
	serviceStatus := 0 // HTTP服务状态，0代表服务未启动，1代表服务已启动（NOTE: 不能在contrilButton按钮内部定义）
	qrStatus := 0      // 二维码状态，0代表二维码未显示，1代表二维码已显示

	// 服务启动/停止按钮逻辑
	controlButton = widget.NewButton("Start", func() {
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
		serviceUrl = fmt.Sprintf("http://%s:%s", selectedInterfaceIP, selectedPort)

		// 生成二维码
		qrCodeImage, err := general.QrCodeImage(serviceUrl)
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
			httpServer, err = HttpServer(selectedInterfaceIP, selectedPort, selectedDir)
			if err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			} else {
				// 设置服务状态
				serviceStatus = 1             // 服务已启动
				statusAnimation.Start()       // 服务状态动画
				controlButton.SetText("Stop") // 修改按钮文字
				// 设置二维码状态
				qrWindow.SetContent(qrImage)                // 将二维码图像添加到窗口（NOTE: 不能使用container.NewCenter()函数将其添加到窗口中心，否则会产生内边距）
				qrWindow.SetPadded(false)                   // 设置窗口内边距为零以确保图像与窗口边框贴合
				qrWindow.Show()                             // 显示二维码窗口
				qrButton.Enable()                           // 启用二维码显示/隐藏按钮
				qrButton.SetIcon(theme.VisibilityOffIcon()) // 按钮变为点击隐藏
				qrStatus = 1                                // 二维码已显示
				// 设置URL按钮
				urlButton.Enable() // 启用URL按钮
				fmt.Printf("\x1b[32;1mServing HTTP on %s port %s (%s)\x1b[0m\n", selectedInterfaceIP, selectedPort, serviceUrl)
			}
		} else if serviceStatus == 1 {
			// 停止HTTP服务
			if err := httpServer.Shutdown(nil); err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			}
			// 设置服务状态
			serviceStatus = 0              // 服务已停止
			statusAnimation.Stop()         // 服务状态动画
			controlButton.SetText("Start") // 修改按钮文字
			// 设置二维码状态
			qrWindow.Hide()                          // 隐藏二维码窗口（NOTE: 不能使用Close()）
			qrButton.Disable()                       // 禁用二维码显示/隐藏按钮
			qrButton.SetIcon(theme.VisibilityIcon()) // 按钮变为点击显示
			qrStatus = 0                             // 二维码未显示
			// 设置URL按钮
			urlButton.Disable() // 禁用URL按钮
		} else {
			customErrText := "Unknown error"
			errorDialog := makeErrorDialog("Error", "Close", customErrText, errorDialogSize, mainWindow)
			errorDialog.Show()
		}
	})
	// 设置按钮外观
	controlButton.Importance = widget.HighImportance // 按钮突出程度

	// 二维码显示/隐藏按钮逻辑
	qrButton = widget.NewButtonWithIcon("", theme.VisibilityIcon(), func() {
		if qrStatus == 0 && serviceStatus == 1 { // 二维码未显示但服务已启动，则显示二维码
			qrWindow.Show()
			qrButton.SetIcon(theme.VisibilityOffIcon()) // 按钮变为点击隐藏
			qrStatus = 1
		} else if qrStatus == 1 && serviceStatus == 1 { // 二维码已显示且服务已启动，则隐藏二维码
			qrWindow.Hide()
			qrButton.SetIcon(theme.VisibilityIcon()) // 按钮变为点击显示
			qrStatus = 0
		}
	})
	qrButton.Disable()                            // 禁用二维码显示/隐藏按钮
	qrButton.Importance = widget.MediumImportance // 按钮突出程度

	// 多态行 —— 服务路径选择按钮 + 已选路径显示框
	crossDirRow := container.NewBorder(nil, nil, folderButton, nil, selectedFolderEntry)
	// 多态行 —— 二维码显示/隐藏按钮 + 服务链接打开按钮 + 状态动画
	crossStatusRow := container.NewBorder(nil, nil, qrButton, urlButton, statusAnimation)

	// 填充主窗口
	content := container.NewVBox(
		interfaceLabel, // 网卡标签
		interfaceRadio, // 网卡选择
		portEntry,      // 端口配置
		crossDirRow,    // 多态行
		separator,      // 分隔线
		crossStatusRow, // 多态行
		separator,      // 分隔线
		controlButton,  // 启动按钮
	)
	mainWindow.SetContent(content)

	// 启动主窗口
	mainWindow.ShowAndRun()
}

// 生成自定义错误提示框
func makeErrorDialog(title, dismiss, text string, size fyne.Size, parent fyne.Window) *dialog.CustomDialog {
	content := widget.NewLabel(text)     // 设置内容
	content.Wrapping = fyne.TextWrapWord // 设置换行
	errorDialog := dialog.NewCustom(title, dismiss, content, parent)
	errorDialog.Resize(size)
	return errorDialog
}

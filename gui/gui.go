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
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/yhyj/skynet/general"
)

// 启动GUI
func StartGraphicalUserInterface() {
	// 获取当前用户信息
	currentUserInfo, err := general.GetCurrentUserInfo()
	if err != nil {
		log.Printf("\x1b[31m%s\x1b[0m\n", err)
	}

	// HTTP服务默认配置
	var (
		defaultIP    = "0.0.0.0"                                           // HTTP服务默认绑定的IP
		defaultPort  = "8080"                                              // HTTP服务默认监听的端口
		defaultDir   = filepath.Join(currentUserInfo.HomeDir, "Downloads") // HTTP服务默认启动路径
		serviceUrl   = fmt.Sprintf("http://%s:%s", defaultIP, defaultPort) // HTTP服务默认URL
		serviceSlice = []string{"Download", "Upload", "All"}               // HTTP服务默认支持启用的方法
	)

	// 界面显示配置
	var (
		serviceLabelText   = "Select Service:"                                                                                  // 服务选择标签默认文本
		interfaceLabelText = "Select Interface:"                                                                                // 网卡选择标签默认文本
		portText           = fmt.Sprintf("Port [1~65535], default %s", defaultPort)                                             // 端口框默认文本
		selectedDirText    = fmt.Sprintf("Directory, default %s", strings.Replace(defaultDir, currentUserInfo.HomeDir, "~", 1)) // 服务启动路径框默认文本
	)

	// 定义服务接口和小部件
	var (
		httpServer    *http.Server    // HTTP服务
		qrWindow      fyne.Window     // 二维码窗口
		windowContent *fyne.Container // 窗口内容容器
		refreshButton *widget.Button  // 接口刷新按钮
		folderButton  *widget.Button  // 目录选择按钮
		qrButton      *widget.Button  // 二维码显示/隐藏按钮
		urlButton     *widget.Button  // 打开URL按钮
		controlButton *widget.Button  // 服务的启动/停止按钮
	)

	// 定义标志位
	var (
		serviceStatus = 0 // HTTP服务状态，0代表服务未启动，1代表服务已启动
		qrStatus      = 0 // 二维码状态，0代表二维码未显示，1代表二维码已显示
	)

	// 定义通用资源
	var (
		separator = widget.NewSeparator() // 创建分隔线
		spacer    = layout.NewSpacer()    // 创建填充空白
	)

	// 创建一个新应用
	appInstance := app.NewWithID(general.Name)
	appInstance.SetIcon(fyne.NewStaticResource("icon", resourceIconPng.StaticContent))

	// 创建主窗口
	mainWindow := appInstance.NewWindow(fmt.Sprintf("%s - %s", general.Name, general.Version))
	mainWindow.SetMaster()                                                                         // 该窗口设为主窗口
	mainWindow.SetFixedSize(false)                                                                 // 是否固定窗口大小
	baseWeight, baseHeight := float32(len(selectedDirText))*9.1, mainWindow.Canvas().Size().Height // 窗口基础尺寸
	mainWindow.Resize(fyne.NewSize(baseWeight, baseHeight))                                        // 设置窗口大小

	// 创建错误提示框尺寸
	errorDialogSize := fyne.NewSize(baseWeight-float32(20), baseHeight-float32(20))

	// 获取网卡信息
	interfaceLabel := widget.NewLabel(interfaceLabelText)
	nicInfos, err := GetNetInterfaces()
	if err != nil {
		errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
		errorDialog.Show()
	}
	// 创建接口选择器（单选按钮组）
	interfaceRadio := widget.NewRadioGroup(nicInfos, func(selected string) {})
	// 创建接口刷新按钮
	refreshButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		nicInfos, err := GetNetInterfaces()
		if err != nil {
			errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
			errorDialog.Show()
		}

		log.Printf("Network interface refresh")
		interfaceRadio.Options = nicInfos
		windowContent.Refresh()
	})

	// 创建端口选择器
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder(portText)
	portEntry.Validator = func(text string) error {
		value, err := strconv.Atoi(text)
		if err != nil || value < 1 || value > 65535 {
			return fmt.Errorf("Invalid port\n")
		}
		return nil
	}

	// 创建目录选择器标签
	selectedDirEntry := widget.NewEntry()
	selectedDirEntry.SetPlaceHolder(selectedDirText)
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
		fileDialog := dialog.NewFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			} else if dir == nil {
				// 未选择文件夹，使用默认值
				selectedDirEntry.SetText(defaultDir)
			} else {
				// 在标签中显示选择的文件夹路径
				selectedDirEntry.SetText(strings.Split(dir.String(), "//")[1])
			}
			// 关闭新窗口
			folderWindow.Close()
		}, folderWindow)
		fileDialog.Show()
		fileDialog.Resize(fyne.NewSize(folderWidth, folderHeight))
	})

	// 创建服务选择标签
	serviceSelectLabel := widget.NewLabel(serviceLabelText)
	// 创建服务选择器
	serviceSelect := widget.NewSelect(serviceSlice, func(selected string) {})
	serviceSelect.Selected = serviceSlice[0]

	// 创建URL打开按钮
	urlButton = widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		serviceUrlParsed, err := url.Parse(serviceUrl)
		if err != nil {
			errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
			errorDialog.Show()
		}
		appInstance.OpenURL(serviceUrlParsed)
		log.Printf("Open URL: \x1b[34;1;4m%s\x1b[0m", serviceUrlParsed)
	})
	urlButton.Disable() // 禁用URL按钮

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

	// 服务启动/停止按钮逻辑
	controlButton = widget.NewButton("Start", func() {
		// 获取参数信息，如果参数为空则使用默认值
		selectedService := func() string {
			if serviceSelect.Selected != "" {
				return serviceSelect.Selected
			}
			return serviceSlice[0]
		}()
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
			if selectedDirEntry.Text != "" {
				return selectedDirEntry.Text
			}
			return defaultDir
		}()

		serviceUrl = func() string {
			if selectedService == "Upload" {
				return fmt.Sprintf("http://%s:%s/upload", selectedInterfaceIP, selectedPort)
			}
			return fmt.Sprintf("http://%s:%s", selectedInterfaceIP, selectedPort)
		}()

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

		if serviceStatus == 0 { // Start
			// 启动HTTP服务
			switch selectedService {
			case "Download":
				httpServer, err = HttpDownloadServer(selectedInterfaceIP, selectedPort, selectedDir)
			case "Upload":
				httpServer, err = HttpUploadServer(selectedInterfaceIP, selectedPort, selectedDir)
			case "All":
				httpServer, err = HttpDownloadUploadServer(selectedInterfaceIP, selectedPort, selectedDir)
			default:
				errorDialog := makeErrorDialog("Warning", "Close", "Please select service", errorDialogSize, mainWindow)
				errorDialog.Show()
			}
			if err != nil {
				errorDialog := makeErrorDialog("Error", "Close", err.Error(), errorDialogSize, mainWindow)
				errorDialog.Show()
			} else {
				// 设置服务状态
				serviceStatus = 1             // 服务已启动
				statusAnimation.Start()       // 服务状态动画
				controlButton.SetText("Stop") // 修改按钮文字
				log.Printf("\x1b[32;1mServing HTTP at '%s'\x1b[0m (\x1b[34;1;4m%s\x1b[0m)\n", selectedDir, serviceUrl)
				// 设置二维码状态
				qrWindow.SetContent(qrImage)                // 将二维码图像添加到窗口（NOTE: 不能使用container.NewCenter()函数将其添加到窗口中心，否则会产生内边距）
				qrWindow.SetPadded(false)                   // 设置窗口内边距为零以确保图像与窗口边框贴合
				qrWindow.Show()                             // 显示二维码窗口
				qrButton.Enable()                           // 启用二维码显示/隐藏按钮
				qrButton.SetIcon(theme.VisibilityOffIcon()) // 变更按钮图标
				qrStatus = 1                                // 二维码已显示
				// 设置URL按钮
				urlButton.Enable() // 启用URL按钮
				// 以下部件禁用
				serviceSelect.Disable()    // 服务选择器
				interfaceRadio.Disable()   // 网卡选择器
				portEntry.Disable()        // 端口输入框
				selectedDirEntry.Disable() // 目录输入框
				folderButton.Disable()     // 目录选择按钮
			}
		} else if serviceStatus == 1 { // Stop
			// 注销处理器
			DeregisterAll(serviceSlice)
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
			qrButton.SetIcon(theme.VisibilityIcon()) // 变更按钮图标
			qrStatus = 0                             // 二维码未显示
			// 设置URL按钮
			urlButton.Disable() // 禁用URL按钮
			// 以下部件启用
			serviceSelect.Enable()    // 服务选择器
			interfaceRadio.Enable()   // 网卡选择器
			portEntry.Enable()        // 端口输入框
			selectedDirEntry.Enable() // 目录输入框
			folderButton.Enable()     // 目录选择按钮
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
			log.Printf("QR Code displayed")
		} else if qrStatus == 1 && serviceStatus == 1 { // 二维码已显示且服务已启动，则隐藏二维码
			qrWindow.Hide()
			qrButton.SetIcon(theme.VisibilityIcon()) // 按钮变为点击显示
			qrStatus = 0
			log.Printf("QR Code hidden")
		}
	})
	qrButton.Disable()                            // 禁用二维码显示/隐藏按钮
	qrButton.Importance = widget.MediumImportance // 按钮突出程度

	// 多态行 —— 服务选择标签 + 服务选择器
	crossServiceRow := container.NewBorder(nil, nil, serviceSelectLabel, nil, serviceSelect)
	// 多态行 —— 接口选择标签 + 接口刷新按钮
	crossInterfaceRow := container.NewBorder(nil, nil, interfaceLabel, refreshButton, nil)
	// 多态行 —— 服务路径选择按钮 + 已选路径显示框
	crossDirRow := container.NewBorder(nil, nil, folderButton, nil, selectedDirEntry)
	// 多态行 —— 二维码显示/隐藏按钮 + 服务链接打开按钮 + 状态动画
	crossStatusRow := container.NewBorder(nil, nil, qrButton, urlButton, statusAnimation)

	// 填充主窗口
	windowContent = container.NewVBox(
		crossServiceRow,   // 多态行 —— 服务选择标签 + 服务选择器
		crossInterfaceRow, // 多态行 —— 接口选择标签 + 接口刷新按钮
		interfaceRadio,    // 接口选择
		spacer,            // 填充空白
		portEntry,         // 端口配置
		crossDirRow,       // 多态行 —— 服务路径选择按钮 + 已选路径显示框
		separator,         // 分隔线
		crossStatusRow,    // 多态行 —— 二维码显示/隐藏按钮 + 服务链接打开按钮 + 状态动画
		separator,         // 分隔线
		controlButton,     // 启动按钮
	)
	mainWindow.SetContent(windowContent)

	// 启动主窗口
	mainWindow.ShowAndRun()
}

// 生成自定义错误提示框
func makeErrorDialog(title, dismiss, text string, size fyne.Size, parent fyne.Window) *dialog.CustomDialog {
	dialogContent := widget.NewLabel(text)     // 设置提示框内容
	dialogContent.Wrapping = fyne.TextWrapWord // 设置换行方式
	errorDialog := dialog.NewCustom(title, dismiss, dialogContent, parent)
	errorDialog.Resize(size)
	return errorDialog
}

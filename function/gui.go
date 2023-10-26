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

func StartGraphicalUserInterface() {
	// HTTP服务变量
	var server *http.Server
	// 参数
	var (
		// HTTP服务启动参数
		selectedInterfaceIP = "0.0.0.0"
		selectedPort        = "8080"
		selectedDir         = "."
	)

	// 创建一个新应用
	app := app.New()
	mainWindow := app.NewWindow("HTTP服务启动器")
	mainWindow.SetMaster() // 设置主窗口
	mainWindow.Resize(fyne.NewSize(300, mainWindow.Canvas().Size().Height))
	mainWindow.SetFixedSize(true) // 固定窗口大小

	// 获取网卡信息
	interfaceLabel := widget.NewLabel("选择网卡:")
	nics, err := GetNetInterfacesForGui()
	if err != nil {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
	}
	// 创建一个单选按钮组
	interfaceRadio := widget.NewRadioGroup(nics, func(selected string) {})

	// 端口选择
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("输入端口号")

	// 目录选择
	selectedFolderEntry := widget.NewEntry()
	selectedFolderEntry.SetPlaceHolder("选择服务目录")
	folderButton := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		// 固定文件选择对话框大小不可修改
		const width float32 = 770
		const height float32 = 481
		// 创建一个新窗口用于文件夹选择
		folderWindow := app.NewWindow("选择服务目录")
		folderWindow.Resize(fyne.NewSize(width, height))
		folderWindow.SetFixedSize(true) // 固定窗口大小
		folderWindow.CenterOnScreen()   // 居中显示

		// 显示新窗口
		folderWindow.Show()

		// 弹出文件夹选择对话框
		fileDialog := dialog.NewFolderOpen(func(dri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
			} else if dri == nil {
				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "Startup folder not set")
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

	// 服务状态显示
	statusAnimation := widget.NewProgressBarInfinite()
	statusAnimation.Stop()

	// 启动按钮
	startButtonClicked := false // 创建一个标志，用于跟踪按钮是否已经被点击
	startButton := widget.NewButton("Start", func() {
		// 检查按钮是否已经被点击，如果已经点击则返回
		if startButtonClicked {
			return
		}
		// 获取参数信息，如果参数为空则使用默认值
		selectedInterfaceIP = func() string {
			parts := strings.Split(interfaceRadio.Selected, " ")
			if len(parts) > 1 {
				return parts[len(parts)-1]
			}
			return "0.0.0.0"
		}()
		selectedPort = func() string {
			if portEntry.Text == "" {
				return "8080"
			}
			return portEntry.Text
		}()
		selectedDir = func() string {
			if selectedFolderEntry.Text == "" {
				return "."
			}
			return selectedFolderEntry.Text
		}()

		url := fmt.Sprintf("http://%s:%v", selectedInterfaceIP, selectedPort)
		fmt.Printf("\x1b[32;1mServing HTTP on %s port %v (%s).\x1b[0m\n", selectedInterfaceIP, selectedPort, url)

		// 启动HTTP服务
		server = HttpServerForGui(selectedInterfaceIP, selectedPort, selectedDir) // 启动HTTP服务
		statusAnimation.Start()                                                   // 启动状态动画

		// 生成二维码
		qrCodeImage, err := QrCodeImage(url)
		if err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		}
		// 将二维码图像转换为 Fyne 图像
		qrImage := canvas.NewImageFromImage(qrCodeImage)
		// 设置图像填充模式为 ImageFillOriginal，以确保不拉伸
		qrImage.FillMode = canvas.ImageFillOriginal

		// 创建一个新窗口用于显示二维码
		qrWindow := app.NewWindow("QR Code")
		// 将二维码图像添加到窗口（
		qrWindow.SetContent(qrImage) // NOTE: 不能使用container.NewCenter()函数将其添加到窗口中心，否则会产生内边距
		// 设置窗口内边距为零以确保图像与窗口边框贴合
		qrWindow.SetPadded(false)
		// 显示窗口
		qrWindow.Show()

		// 按钮已被点击
		startButtonClicked = true
	})

	// 停止按钮
	stopButton := widget.NewButton("Stop", func() {
		fmt.Println("Stop service")

		// 在这里执行停止HTTP服务的代码
		if err := server.Shutdown(nil); err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		}
		statusAnimation.Stop()
	})

	content := container.NewVBox(
		interfaceLabel, interfaceRadio, // 网卡选择
		portEntry,       // 端口配置
		dirRow,          // 文件夹选择
		statusAnimation, // 状态显示
		startButton,     // 启动按钮
		stopButton,      // 停止按钮
	)
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
}

/*
File: gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 08:52:06

Description: 执行子命令 'gui'
*/

package cmd

import (
	"log"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/skynet/general"
	"github.com/yhyj/skynet/gui"
)

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start the GUI version of skynet",
	Long:  `Start the skynet Graphical User Interface`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			errorSlogan []string // 错误标语
		)

		// 检查平台
		if general.Platform == "linux" {
			// 检查是否远程连接
			if general.GetVariable("DISPLAY") != "" {
				// 设置字体
				if err := gui.SetFont(); err != nil {
					log.Println(general.FgRedText(err))
				}
				// 启动 GUI
				gui.StartGraphicalUserInterface()
			} else {
				errorSlogan = append(errorSlogan, "Could not connect to display")
			}
		} else if general.Platform == "windows" {
			// 设置字体
			if err := gui.SetFont(); err != nil {
				log.Println(general.FgRedText(err))
			}
			// 启动 GUI
			gui.StartGraphicalUserInterface()
		} else if general.Platform == "darwin" {
			// 设置字体
			if err := gui.SetFont(); err != nil {
				log.Println(general.FgRedText(err))
			}
			// 启动 GUI
			gui.StartGraphicalUserInterface()
		} else {
			errorSlogan = append(errorSlogan, "Current platform is not supported")
		}

		// 输出标语
		if len(errorSlogan) > 0 {
			color.Println()
			for _, slogan := range errorSlogan {
				color.Error.Tips(general.PrimaryText(slogan))
			}
		}
	},
}

func init() {
	guiCmd.Flags().BoolP("help", "h", false, "help for gui command")
	rootCmd.AddCommand(guiCmd)
}

/*
File: gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 08:52:06

Description: 由程序子命令 gui 执行
*/

package cmd

import (
	"fmt"
	"log"

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
		if general.Platform == "linux" {
			if general.GetVariable("DISPLAY") != "" {
				// 设置字体
				if err := gui.SetFont(); err != nil {
					log.Printf(general.ErrorBaseFormat, err)
				}
				// 启动 GUI
				gui.StartGraphicalUserInterface()
			} else {
				fmt.Println("Could not connect to display")
			}
		} else if general.Platform == "windows" {
			// 设置字体
			if err := gui.SetFont(); err != nil {
				log.Printf(general.ErrorBaseFormat, err)
			}
			// 启动 GUI
			gui.StartGraphicalUserInterface()
		} else if general.Platform == "darwin" {
			// 设置字体
			if err := gui.SetFont(); err != nil {
				log.Printf(general.ErrorBaseFormat, err)
			}
			// 启动 GUI
			gui.StartGraphicalUserInterface()
		} else {
			fmt.Println("Current platform is not supported")
		}
	},
}

func init() {
	guiCmd.Flags().BoolP("help", "h", false, "help for gui command")
	rootCmd.AddCommand(guiCmd)
}

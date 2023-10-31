/*
File: gui.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-10-26 08:52:06

Description: 程序子命令'gui'时执行
*/

package cmd

import (
	"fmt"

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
					fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
				}
				// 启动GUI
				gui.StartGraphicalUserInterface()
			} else {
				fmt.Println("The DISPLAY environment variable is missing, please use the CLI version")
			}
		} else if general.Platform == "windows" {
			// 设置字体
			if err := gui.SetFont(); err != nil {
				fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
			}
			// 启动GUI
			gui.StartGraphicalUserInterface()
		} else if general.Platform == "darwin" {
			fmt.Println("macOS platform is not supported yet")
		} else {
			fmt.Println("Current platform is not supported")
		}
	},
}

func init() {
	guiCmd.Flags().BoolP("help", "h", false, "help for gui command")
	rootCmd.AddCommand(guiCmd)
}

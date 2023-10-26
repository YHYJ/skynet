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
	"github.com/yhyj/skynet/function"
)

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Start the GUI version of skynet",
	Long:  `Start the skynet Graphical User Interface`,
	Run: func(cmd *cobra.Command, args []string) {
		// 设置字体
		if err := function.SetFont(); err != nil {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", err)
		}
		// 启动GUI
		function.StartGraphicalUserInterface()
	},
}

func init() {
	guiCmd.Flags().BoolP("help", "h", false, "help for httpserver command")
	rootCmd.AddCommand(guiCmd)
}

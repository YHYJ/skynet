/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 10:18:55

Description: 程序未带子命令或参数时执行
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yhyj/skynet/function"
)

var rootCmd = &cobra.Command{
	Use:   "skynet",
	Short: function.Translate(function.Localizer, "CmdRootShort", function.TemplateData),
	Long:  function.Translate(function.Localizer, "CmdRootLong", function.TemplateData),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, function.Translate(function.Localizer, "CmdRootHelpFlag", function.TemplateData))
}

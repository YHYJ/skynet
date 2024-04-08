/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-05-25 10:18:55

Description: 执行程序
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "skynet",
	Short: "Network affairs administrator",
	Long:  `skynet is the system network affairs administrator.`,
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
	rootCmd.Flags().BoolP("help", "h", false, "help for skynet")
}

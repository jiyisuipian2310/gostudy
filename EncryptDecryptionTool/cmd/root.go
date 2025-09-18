/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "EncryptDecryptionTool",
	Short: "AES SM4 加解密工具",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 这个函数会在任何子命令执行前运行
		fmt.Println("I am the ROOT PersistentPreRun hook! I run before every command.\n")
	},
	PersistentPostRun : func(cmd *cobra.Command, args []string) {
		fmt.Println("I am the ROOT PersistentPreRun hook! I run after every command.\n")
	},
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 设置自定义帮助模板

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.EncryptDecryptionTool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

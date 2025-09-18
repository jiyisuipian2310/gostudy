/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sm4DecryptCmd represents the sm4Decrypt command
var sm4DecryptCmd = &cobra.Command{
	Use:   "sm4Decrypt",
	Short: "SM4 解密",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sm4Decrypt called")
	},
}

func init() {
	rootCmd.AddCommand(sm4DecryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sm4DecryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sm4DecryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

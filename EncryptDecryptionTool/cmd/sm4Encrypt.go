/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sm4EncryptCmd represents the sm4Encrypt command
var sm4EncryptCmd = &cobra.Command{
	Use:   "sm4Encrypt",
	Short: "SM4 加密",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sm4Encrypt called")
	},
}

func init() {
	rootCmd.AddCommand(sm4EncryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sm4EncryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sm4EncryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

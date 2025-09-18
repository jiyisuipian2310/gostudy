/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agclqq/goencryption"
)

var bAesEncryptCmdMybool = false

// aesEncryptCmd represents the aesEncrypt command
var aesEncryptCmd = &cobra.Command{
	Use:   "aesEncrypt",
	Short: "AES 加密",
	Args:  cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("aesEncryptCmd hook, just before Run.\n")
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("\naesEncryptCmd hook, just after Run.\n")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bAesEncryptCmdMybool: %t\n", bAesEncryptCmdMybool)
	

		fmt.Printf("请输入AES Key(默认: 63dTjxISXlwAso0n): ")
		
		var strAesKey = ""
		fmt.Scanln(&strAesKey)
		if(strAesKey == "") {
			strAesKey = "63dTjxISXlwAso0n"
		}
		
		fmt.Print("请输入AES IV(默认: a1b2c3d4e5f6g7h8): ")
		var strAesIV = ""
		fmt.Scanln(&strAesIV)
		if(strAesIV == "") {
			strAesIV = "a1b2c3d4e5f6g7h8"
		}

		fmt.Print("请输入明文数据: ")
		var strPlainData = ""
		fmt.Scanln(&strPlainData)
		if(strPlainData == "") {
			fmt.Printf("没有输入明文数据，程序退出！\n")
			return
		}

		cipherData, err := goencryption.EasyEncrypt("aes/cbc/pkcs7/base64", strPlainData, strAesKey, strAesIV)
		if err != nil {
			fmt.Printf("EasyEncrypt Failed: %s\n", err.Error())
			return
		}

		fmt.Printf("CipherData: %s\n", cipherData)
	},
}

func init() {
	rootCmd.AddCommand(aesEncryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aesEncryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	aesEncryptCmd.Flags().BoolVarP(&bAesEncryptCmdMybool, "toggle", "t", false, "This is bool value flags")
}

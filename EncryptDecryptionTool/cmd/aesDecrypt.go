/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/agclqq/goencryption"
)

var bAesDecryptCmdMybool = false

// aesDecryptCmd represents the aesDecrypt command
var aesDecryptCmd = &cobra.Command{
	Use:   "aesDecrypt",
	Short: "AES 解密",
    PreRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("aesDecryptCmd hook, just before Run.\n")
    },
	PostRun: func(cmd *cobra.Command, args []string) {
        fmt.Println("\naesDecryptCmd hook, just after Run.\n")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bAesDecryptCmdMybool: %t\n", bAesDecryptCmdMybool)

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

        fmt.Print("请输入密文数据: ")
        var strCipherData = ""
        fmt.Scanln(&strCipherData)
        if(strCipherData == "") {
            fmt.Printf("没有输入密文数据，程序退出！\n")
            return
        }

        plainData, err := goencryption.EasyDecrypt("aes/cbc/pkcs7/base64", strCipherData, strAesKey, strAesIV)
        if err != nil {
            fmt.Printf("EasyDecrypt Failed: %s\n", err.Error())
            return
        }

        fmt.Printf("plainData: %s\n", plainData)
	},
}

func init() {
	rootCmd.AddCommand(aesDecryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aesDecryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	aesDecryptCmd.Flags().BoolVarP(&bAesDecryptCmdMybool, "toggle", "t", false, "This is bool value flags")
}

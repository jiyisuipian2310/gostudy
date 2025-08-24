package interface_libs_src

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"unsafe"
)

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

// 填充数据到块大小
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// 去除填充数据
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

//export AESCBCEncrypt
func AESCBCEncrypt(plaintext *C.char, key *C.char, iv *C.char, errMsg *C.char) *C.char {
	// 初始化错误消息
	if errMsg != nil {
		C.memset(unsafe.Pointer(errMsg), 0, 256)
	}

	// 将C字符串转换为Go字符串
	goPlaintext := C.GoString(plaintext)
	goKey := C.GoString(key)
	goIV := C.GoString(iv)

	// 检查key长度，AES需要16, 24, 或32字节的key
	if len(goKey) != 16 && len(goKey) != 24 && len(goKey) != 32 {
		if errMsg != nil {
			errorMsg := "Key must be 16, 24, or 32 bytes long, got " + string(rune(len(goKey))) + " bytes"
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 检查IV长度，必须是16字节
	if len(goIV) != 16 {
		if errMsg != nil {
			errorMsg := "IV must be 16 bytes long, got " + string(rune(len(goIV))) + " bytes"
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 创建AES cipher
	block, err := aes.NewCipher([]byte(goKey))
	if err != nil {
		if errMsg != nil {
			errorMsg := "Error creating cipher: " + err.Error()
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 对明文进行PKCS7填充
	plaintextBytes := []byte(goPlaintext)
	plaintextBytes = pkcs7Padding(plaintextBytes, aes.BlockSize)

	// 创建CBC模式加密器
	mode := cipher.NewCBCEncrypter(block, []byte(goIV))

	// 加密
	ciphertext := make([]byte, len(plaintextBytes))
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// Base64编码
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return C.CString(encoded)
}

//export AESCBCDecrypt
func AESCBCDecrypt(ciphertext *C.char, key *C.char, iv *C.char, errMsg *C.char) *C.char {
	// 初始化错误消息
	if errMsg != nil {
		C.memset(unsafe.Pointer(errMsg), 0, 256)
	}

	// 将C字符串转换为Go字符串
	goCiphertext := C.GoString(ciphertext)
	goKey := C.GoString(key)
	goIV := C.GoString(iv)

	// 检查key长度
	if len(goKey) != 16 && len(goKey) != 24 && len(goKey) != 32 {
		if errMsg != nil {
			errorMsg := "Key must be 16, 24, or 32 bytes long, got " + string(rune(len(goKey))) + " bytes"
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 检查IV长度
	if len(goIV) != 16 {
		if errMsg != nil {
			errorMsg := "IV must be 16 bytes long, got " + string(rune(len(goIV))) + " bytes"
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// Base64解码
	decoded, err := base64.StdEncoding.DecodeString(goCiphertext)
	if err != nil {
		if errMsg != nil {
			errorMsg := "Error decoding base64: " + err.Error()
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 创建AES cipher
	block, err := aes.NewCipher([]byte(goKey))
	if err != nil {
		if errMsg != nil {
			errorMsg := "Error creating cipher: " + err.Error()
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 检查密文长度是否为块大小的倍数
	if len(decoded)%aes.BlockSize != 0 {
		if errMsg != nil {
			errorMsg := "Ciphertext is not a multiple of the block size"
			C.strncpy(errMsg, C.CString(errorMsg), 255)
		}
		return nil
	}

	// 创建CBC模式解密器
	mode := cipher.NewCBCDecrypter(block, []byte(goIV))

	// 解密
	plaintext := make([]byte, len(decoded))
	mode.CryptBlocks(plaintext, decoded)

	// 去除PKCS7填充
	plaintext = pkcs7UnPadding(plaintext)

	return C.CString(string(plaintext))
}

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

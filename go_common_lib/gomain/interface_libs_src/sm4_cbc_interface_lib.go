package interface_libs_src

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
    "unsafe"
    "crypto/cipher"
    "fmt"
    "github.com/tjfoc/gmsm/sm4"
    "encoding/base64"
)

// SM4加密 (Go内部使用)
func encryptSM4(key, iv, plaintext []byte) ([]byte, error) {
    block, err := sm4.NewCipher(key)
    if err != nil {
        return nil, err
    }

    // 使用CBC模式
    blockSize := block.BlockSize()

	if len(iv) != blockSize {
		return nil, fmt.Errorf("IV length must equal block size(%d)", blockSize)
	}

    plaintext = pkcs7Padding(plaintext, blockSize)
    ciphertext := make([]byte, len(plaintext))

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(ciphertext, plaintext)

    return ciphertext, nil
}

// SM4解密 (Go内部使用)
func decryptSM4(key, iv, ciphertext []byte) ([]byte, error) {
    block, err := sm4.NewCipher(key)
    if err != nil {
        return nil, err
    }

    // 使用CBC模式
    blockSize := block.BlockSize()
    if len(ciphertext) < blockSize || len(ciphertext)%blockSize != 0 {
        return nil, fmt.Errorf("ciphertext length is not a multiple of the block size")
    }

	if len(iv) != blockSize {
		return nil, fmt.Errorf("IV length must equal block size(%d)", blockSize)
	}

    plaintext := make([]byte, len(ciphertext))

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(plaintext, ciphertext)
    plaintext = pkcs7UnPadding(plaintext)

    return plaintext, nil
}

// 导出函数: SM4加密
//export SM4CBCEncrypt
func SM4CBCEncrypt(plaintext *C.char, key *C.char, iv *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
    if errMsg == nil || errMsgLen < 128 {
        return nil
    }

    C.memset(unsafe.Pointer(errMsg), 0, C.size_t(errMsgLen))

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

    goKey := C.GoString(key)
    goPlaintext := C.GoString(plaintext)
    goIV := C.GoString(iv)

    // 调用加密函数
    ciphertext, err := encryptSM4([]byte(goKey), []byte(goIV), []byte(goPlaintext))
    if err != nil {
        errorMsg := "encryptSM4 failed: " + err.Error()
        setErrorMessage(errorMsg)
        return nil
    }

    encoded := base64.StdEncoding.EncodeToString(ciphertext)

    // 将结果转换回C字符串
    return C.CString(encoded)
}

// 导出函数: SM4解密
//export SM4CBCDecrypt
func SM4CBCDecrypt(ciphertext *C.char, key *C.char, iv *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
    if errMsg == nil || errMsgLen < 128 {
        return nil
    }

    C.memset(unsafe.Pointer(errMsg), 0, C.size_t(errMsgLen))

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

    goKey := C.GoString(key)
    goCiphertext := C.GoString(ciphertext)
    goIV := C.GoString(iv)

    decoded, err := base64.StdEncoding.DecodeString(goCiphertext)
    if err != nil {
        errorMsg := "Error decoding base64: " + err.Error()
        setErrorMessage(errorMsg)
        return nil
    }

    // 调用解密函数
    plaintext, err := decryptSM4([]byte(goKey), []byte(goIV), decoded)
    if err != nil {
        errorMsg := "decryptSM4 failed: " + err.Error()
        setErrorMessage(errorMsg)
        return nil
    }

    // 将结果转换回C字符串
    return C.CString(string(plaintext))
}

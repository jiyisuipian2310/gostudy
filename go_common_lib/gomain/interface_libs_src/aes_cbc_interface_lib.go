package interface_libs_src

import (
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

//export AESCBCEncrypt
func AESCBCEncrypt(plaintext *C.char, key *C.char, iv *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
    if errMsg == nil || errMsgLen < 128 {
        return nil
    }

    // 初始化错误消息
    C.memset(unsafe.Pointer(errMsg), 0, C.size_t(errMsgLen))

    // 将C字符串转换为Go字符串
    goPlaintext := C.GoString(plaintext)
    goKey := C.GoString(key)
    goIV := C.GoString(iv)

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

    // 检查key长度，AES需要16, 24, 或32字节的key
    if len(goKey) != 16 && len(goKey) != 24 && len(goKey) != 32 {
        errorMsg := "Key must be 16, 24, or 32 bytes long, got " + string(rune(len(goKey))) + " bytes"
        setErrorMessage(errorMsg)
        return nil
    }

    // 检查IV长度，必须是16字节
    if len(goIV) != 16 {
        errorMsg := "IV must be 16 bytes long, got " + string(rune(len(goIV))) + " bytes"
        setErrorMessage(errorMsg)
        return nil
    }

    // 创建AES cipher
    block, err := aes.NewCipher([]byte(goKey))
    if err != nil {
        errorMsg := "Error creating cipher: " + err.Error()
        setErrorMessage(errorMsg)
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
func AESCBCDecrypt(ciphertext *C.char, key *C.char, iv *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
    if errMsg == nil || errMsgLen < 128 {
        return nil
    }

    C.memset(unsafe.Pointer(errMsg), 0, C.size_t(errMsgLen))

    // 将C字符串转换为Go字符串
    goCiphertext := C.GoString(ciphertext)
    goKey := C.GoString(key)
    goIV := C.GoString(iv)

    setErrorMessage := func(message string) {
        cMsg := C.CString(message)
        defer C.free(unsafe.Pointer(cMsg)) // 确保释放内存
        C.strncpy(errMsg, cMsg, C.size_t(errMsgLen)-1)
    }

    // 检查key长度
    if len(goKey) != 16 && len(goKey) != 24 && len(goKey) != 32 {
        errorMsg := "Key must be 16, 24, or 32 bytes long, got " + string(rune(len(goKey))) + " bytes"
        setErrorMessage(errorMsg)
        return nil
    }

    // 检查IV长度
    if len(goIV) != 16 {
        errorMsg := "IV must be 16 bytes long, got " + string(rune(len(goIV))) + " bytes"
        setErrorMessage(errorMsg)
        return nil
    }

    // Base64解码
    decoded, err := base64.StdEncoding.DecodeString(goCiphertext)
    if err != nil {
        errorMsg := "Error decoding base64: " + err.Error()
        setErrorMessage(errorMsg)
        return nil
    }

    // 创建AES cipher
    block, err := aes.NewCipher([]byte(goKey))
    if err != nil {
        errorMsg := "Error creating cipher: " + err.Error()
        setErrorMessage(errorMsg)
        return nil
    }

    // 检查密文长度是否为块大小的倍数
    if len(decoded)%aes.BlockSize != 0 {
        errorMsg := "Ciphertext is not a multiple of the block size"
        setErrorMessage(errorMsg)
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

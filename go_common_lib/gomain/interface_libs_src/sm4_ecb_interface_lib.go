package interface_libs_src

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
    "unsafe"
    "fmt"
    "github.com/tjfoc/gmsm/sm4"
    "encoding/base64"
)

// SM4-ECB加密
func encryptSM4ECB(key, plaintext []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create SM4 cipher: %v", err)
	}

	// ECB模式不需要IV
	blockSize := block.BlockSize()
	
	// PKCS7填充
	paddedPlaintext := pkcs7Padding(plaintext, blockSize)
	ciphertext := make([]byte, len(paddedPlaintext))

	// ECB模式加密：逐块加密
	for i := 0; i < len(paddedPlaintext); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], paddedPlaintext[i:i+blockSize])
	}

	return ciphertext, nil
}

// SM4-ECB解密
func decryptSM4ECB(key, ciphertext []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create SM4 cipher: %v", err)
	}

	blockSize := block.BlockSize()
	
	// 检查密文长度
	if len(ciphertext) == 0 {
		return nil, fmt.Errorf("ciphertext is empty")
	}
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext length is not a multiple of block size (%d)", blockSize)
	}

	plaintext := make([]byte, len(ciphertext))

	// ECB模式解密：逐块解密
	for i := 0; i < len(ciphertext); i += blockSize {
		block.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	// 去除PKCS7填充
	return pkcs7UnPadding(plaintext), nil
}

//export SM4ECBEncrypt
func SM4ECBEncrypt(plaindata *C.char, key *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
	if errMsg == nil || errMsgLen < 128 {
		return nil
	}

	setErrorMessage := func(message string) {
		cMsg := C.CString(message)
		defer C.free(unsafe.Pointer(cMsg))
		C.strncpy(errMsg, cMsg, C.size_t(errMsgLen-1))
	}

	if plaindata == nil || key == nil {
		setErrorMessage("SM4ECBEncrypt failed: plaindata or key is NULL")
		return nil
	}

	goPlaintext := C.GoString(plaindata)
	goKey := C.GoString(key)

	if len(goKey) != 16 {
		setErrorMessage("SM4ECBEncrypt failed: SM4 key must be 16 bytes")
		return nil
	}

	// 执行加密
	ciphertext, err := encryptSM4ECB([]byte(goKey), []byte(goPlaintext))
	if err != nil {
		setErrorMessage("SM4ECBEncrypt failed: Encryption failed: " + err.Error())
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return C.CString(encoded)
}

//export SM4ECBDecrypt
func SM4ECBDecrypt(cipherdata *C.char, key *C.char, errMsg *C.char, errMsgLen C.int) *C.char {
	if errMsg == nil || errMsgLen < 128 {
		return nil
	}

	setErrorMessage := func(message string) {
		cMsg := C.CString(message)
		defer C.free(unsafe.Pointer(cMsg))
		C.strncpy(errMsg, cMsg, C.size_t(errMsgLen-1))
	}

	if cipherdata == nil || key == nil {
		setErrorMessage("SM4ECBDecrypt failed: cipherdata or key is NULL")
		return nil
	}

	goCiphertext := C.GoString(cipherdata)
	goKey := C.GoString(key)

	if len(goKey) != 16 {
		setErrorMessage("SM4ECBDecrypt failed: SM4 key must be 16 bytes")
		return nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(goCiphertext)
	if err != nil {
		setErrorMessage("SM4ECBDecrypt failed: Ciphertext base64 decode failed: " + err.Error())
		return nil
	}

	// 执行解密
	plaintext, err := decryptSM4ECB([]byte(goKey), ciphertext)
	if err != nil {
		setErrorMessage("Decryption failed: " + err.Error())
		return nil
	}

	return C.CString(string(plaintext))
}

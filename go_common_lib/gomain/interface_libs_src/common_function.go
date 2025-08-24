package interface_libs_src

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import "bytes"
import "unsafe"

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

//export FreeCString
func FreeCString(str *C.char) {
	C.free(unsafe.Pointer(str))
}

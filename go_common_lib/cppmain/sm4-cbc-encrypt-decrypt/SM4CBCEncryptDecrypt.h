#ifndef __SM4_CRYPTO_H_
#define __SM4_CRYPTO_H_

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

// SM4加密函数
// key: 16字节密钥
// plaintext: 明文数据
// iv: 初始化向量（16字节的字符串）
// 返回: 密文字符串指针，需要调用FreeCString释放内存
char* SM4Encrypt(const char* plaintext, const char* key, const char* iv, char* errMsg, int errMsgLen);

// SM4解密函数
// ciphertext: 密文数据
// key: 16字节密钥
// iv: 初始化向量（16字节的字符串）
// 返回: 明文字符串指针，需要调用FreeCString释放内存
char* SM4Decrypt(const char* ciphertext, const char* key, const char* iv, char* errMsg, int errMsgLen);

// 释放C字符串内存
void FreeCString(char* str);

#ifdef __cplusplus
}
#endif

#endif // __SM4_CRYPTO_H_

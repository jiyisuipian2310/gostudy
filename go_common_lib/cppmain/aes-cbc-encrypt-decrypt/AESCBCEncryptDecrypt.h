#ifndef __AES_CBC_H_
#define __AES_CBC_H_

#ifdef __cplusplus
extern "C" {
#endif

// 加密函数
// plaintext: 要加密的明文
// key: 加密密钥（16, 24, 或32字节的字符串）
// iv: 初始化向量（16字节的字符串）
// errMsg: 错误信息输出参数（如果为NULL表示成功，需要调用FreeCString释放）
// 返回: base64编码的加密结果，需要调用FreeCString释放内存；如果失败返回NULL
char* AESCBCEncrypt(const char* plaintext, const char* key, const char* iv, char* errMsg, int errMsgLen);

// 解密函数
// ciphertext: base64编码的密文
// key: 解密密钥（16, 24, 或32字节的字符串）
// iv: 初始化向量（16字节的字符串）
// errMsg: 错误信息输出参数（如果为NULL表示成功，需要调用FreeCString释放）
// 返回: 解密后的明文，需要调用FreeCString释放内存；如果失败返回NULL
char* AESCBCDecrypt(const char* ciphertext, const char* key, const char* iv, char* errMsg, int errMsgLen);

// 释放由Go返回的字符串内存
void FreeCString(char* str);

#ifdef __cplusplus
}
#endif

#endif // GO_AES_CBC_H

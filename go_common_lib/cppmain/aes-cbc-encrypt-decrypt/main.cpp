#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include "AESCBCEncryptDecrypt.h"
using namespace std;

void test_encrypt_decrypt() {
	string strKey = "63dTjxISXlwAso0n";
	string strIV = "a1b2c3d4e5f6g7h8";
	cout << "秘钥: " << strKey << endl;
	cout << "IV: " << strIV << endl;
	
	string strPlainData = "the world is so big, I want to see see !";
	cout << "原始明文: " << strPlainData << endl;

	char errMsg[256] = { 0 };
    
    char* encrypted = AESCBCEncrypt(strPlainData.data(), strKey.data(), strIV.data(), errMsg);
	if(encrypted == NULL) {
		cout << "加密错误: " << errMsg << endl;
        return;
    }
    
    printf("加密结果: %s\n", encrypted);
    

    char* decrypted = AESCBCDecrypt(encrypted, strKey.data(), strIV.data(), errMsg);
	if(decrypted == NULL) {
		cout << "解密错误: " << errMsg << endl;
        FreeCString(encrypted);
        return;
    }
    
    printf("解密结果: %s\n", decrypted);
    
    // 验证加解密是否正确
    if (strcmp(strPlainData.data(), decrypted) == 0) {
        printf("✓ 加解密测试成功！\n");
    } else {
        printf("✗ 加解密测试失败！\n");
    }
    
    // 释放内存
    FreeCString(encrypted);
    FreeCString(decrypted);
}

void test_error_cases() {
    printf("\n测试错误情况:\n");
    
    // 测试错误的密钥长度
    const char* plaintext = "test";
    const char* wrong_key = "shortkey12345678"; // 错误的密钥长度
    const char* iv = "abcdefghijklmno";
    
    char errMsg[256] = { 0 };
    char* encrypted = AESCBCEncrypt(plaintext, wrong_key, iv, errMsg);
    if (encrypted == NULL) {
        printf("预期错误: %s\n", errMsg);
    }
    
    if (encrypted != NULL) {
        FreeCString(encrypted);
    }
}

int main() {
    test_encrypt_decrypt();
    test_error_cases();
    return 0;
}

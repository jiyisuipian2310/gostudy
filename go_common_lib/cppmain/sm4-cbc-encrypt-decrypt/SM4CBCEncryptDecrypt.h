#ifndef __SM4_CRYPTO_H_
#define __SM4_CRYPTO_H_

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

char* SM4CBCEncrypt(const char* plaintext, const char* key, const char* iv, char* errMsg, int errMsgLen);

char* SM4CBCDecrypt(const char* ciphertext, const char* key, const char* iv, char* errMsg, int errMsgLen);

void FreeCString(char* str);

#ifdef __cplusplus
}
#endif

#endif // __SM4_CRYPTO_H_

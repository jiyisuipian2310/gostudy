#ifndef __SM4_CRYPTO_H_
#define __SM4_CRYPTO_H_

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

char* UploadFileByPassword(const char* address, const char* username, const char* password, const char* remotePath, const char* localPath);
char* DownloadFileByPassword(const char* address, const char* username, const char* password, const char* remotePath, const char* localPath);
void FreeCString(char* str);

#ifdef __cplusplus
}
#endif

#endif // __SM4_CRYPTO_H_

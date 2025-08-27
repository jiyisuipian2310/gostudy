#ifndef __UPLOAD_DOWNLOAD_FILE_H_
#define __UPLOAD_DOWNLOAD_FILE_H_

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

//通过密码登录上传文件
char* UploadFileByPassword(const char* address, const char* username, const char* password, 
	const char* remotePath, const char* localPath);

//通过密码登录下载文件
char* DownloadFileByPassword(const char* address, const char* username, const char* password, 
	const char* remotePath, const char* localPath);

//通过秘钥登录上传文件
char* UploadFileByKey(const char* address, const char* username, const char* privatekey,
		const char* passphrase, const char* remotePath, const char* localPath);

//通过秘钥登录下载文件
char* DownloadFileByKey(const char* address, const char* username, const char* privatekey,
		const char* passphrase, const char* remotePath, const char* localPath);

void FreeCString(char* str);

#ifdef __cplusplus
}
#endif

#endif // __UPLOAD_DOWNLOAD_FILE_H_

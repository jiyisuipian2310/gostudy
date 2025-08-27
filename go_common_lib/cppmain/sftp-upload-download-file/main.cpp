#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include "UploadDownloadFile.h"
using namespace std;


char* read_file_content(char* pFile, char* pMode)
{
	FILE* fp = fopen(pFile, pMode);
	if(fp == NULL) return NULL;

	fseek(fp, 0, SEEK_END);
	long file_size = ftell(fp);
	fseek(fp, 0, SEEK_SET);
	char* buffer = (char*) malloc(file_size + 1);
	fread(buffer, 1, file_size, fp);
	buffer[file_size] = '\0';
	fclose(fp);
	return buffer;	
}

int main()
{
    string strAddr = "192.168.104.100:22";
	string strUser = "root";
	string strPassword = "root100^YHN";

#if 0
    //密码登录
    #if 0
	    string strRemoteFile = "/home/yull/thrift_server.zip";
	    string strLocalFile = "./thrift_server.zip";

	    char* errMsg = DownloadFileByPassword(strAddr.data(), strUser.data(), strPassword.data(), strRemoteFile.data(), strLocalFile.data());
	    if(errMsg == NULL) {
		    cout << "File downloaded successfully!" << endl;
	    } else {
		    cout << "File downloaded failed, errorMsg: " << errMsg << endl;
		    FreeCString(errMsg);
	    }
    #else
	    string strLocalFile = "/root/test.tar.gz";
	    string strRemoteFile = "/home/yull/test.tar.gz";

	    char* errMsg = UploadFileByPassword((char*)strAddr.data(), (char*)strUser.data(), (char*)strPassword.data(), (char*)strRemoteFile.data(), (char*)strLocalFile.data());
	    if(errMsg == NULL) {
		    cout << "File upload successfully!" << endl;
	    } else {
		    cout << "File uploaded failed, errorMsg: " << errMsg << endl;
		    FreeCString(errMsg);
	    }
	
    #endif
#else
	//秘钥登录
	
	char* buffer = read_file_content((char*)"/root/.ssh/id_rsa", (char*)"rb");
	if(buffer == NULL) { return 0; }
	string privatekey = buffer;
	free(buffer);
	string passphrase = "123456";
	
	#if 0
		string strRemoteFile = "/home/yull/thrift_server.zip";
		string strLocalFile = "./thrift_server.zip";		
		char* errMsg = DownloadFileByKey(strAddr.data(), strUser.data(), privatekey.data(), passphrase.data(), strRemoteFile.data(), strLocalFile.data());
	    if(errMsg == NULL) {
		    cout << "File downloaded successfully!" << endl;
	    } else {
		    cout << "File downloaded failed, errorMsg: " << errMsg << endl;
		    FreeCString(errMsg);
	    }
		
	#else
	    string strLocalFile = "/root/test.tar.gz";
	    string strRemoteFile = "/home/yull/test.tar.gz";
		char* errMsg = UploadFileByKey(strAddr.data(), strUser.data(), privatekey.data(), passphrase.data(), strRemoteFile.data(), strLocalFile.data());
	    if(errMsg == NULL) {
		    cout << "File upload successfully!" << endl;
	    } else {
		    cout << "File uploaded failed, errorMsg: " << errMsg << endl;
		    FreeCString(errMsg);
	    }
	#endif
#endif

    return 0;
}

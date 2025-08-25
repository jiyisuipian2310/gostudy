#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include "UploadDownloadFile.h"
using namespace std;

int main()
{
    string strAddress = "192.168.104.100:22";
	string strUserName = "root";
	string strPassword = "root100^YHN";
	char errorMsg[128] = {0};

#if 0
	string strRemoteFile = "/home/yull/thrift_server.zip";
	string strLocalFile = "./thrift_server.zip";

	int result = DownloadFileByPassword(strAddress.data(), strUserName.data(), strPassword.data(), 
		strRemoteFile.data(), strLocalFile.data(), errorMsg, sizeof(errorMsg));
	if(result == 0) {
		cout << "File downloaded successfully!" << endl;
	} else {
		cout << "File downloaded failed, errorMsg: " << errorMsg << endl;
	}
#else
	string strLocalFile = "/root/test.tar.gz22";
	string strRemoteFile = "/home/yull/test.tar.gz";

	int result = UploadFileByPassword(strAddress.data(), strUserName.data(), strPassword.data(), 
		strRemoteFile.data(), strLocalFile.data(), errorMsg, sizeof(errorMsg));
	if(result == 0) {
		cout << "File upload successfully!" << endl;
	} else {
		cout << "File uploaded failed, errorMsg: " << errorMsg << endl;
	}
	
#endif

    return 0;
}

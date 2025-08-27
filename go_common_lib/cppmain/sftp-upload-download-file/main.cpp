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

#if 0
	string strRemoteFile = "/home/yull/thrift_server.zip";
	string strLocalFile = "./thrift_server.zip";

	char* errMsg = DownloadFileByPassword(strAddress.data(), strUserName.data(), strPassword.data(), strRemoteFile.data(), strLocalFile.data());
	if(errMsg == NULL) {
		cout << "File downloaded successfully!" << endl;
	} else {
		cout << "File downloaded failed, errorMsg: " << errMsg << endl;
		FreeCString(errMsg);
	}
#else
	string strLocalFile = "/root/test.tar.gz";
	string strRemoteFile = "/home/yull/test.tar.gz";

	char* errMsg = UploadFileByPassword(strAddress.data(), strUserName.data(), strPassword.data(), strRemoteFile.data(), strLocalFile.data());
	if(errMsg == NULL) {
		cout << "File upload successfully!" << endl;
	} else {
		cout << "File uploaded failed, errorMsg: " << errMsg << endl;
		FreeCString(errMsg);
	}
	
#endif

    return 0;
}

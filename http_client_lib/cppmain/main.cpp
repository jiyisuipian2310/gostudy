#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <iostream>
#include <string>
#include "IHttpBase.h"
using namespace std;

class HttpClient: public IHttpBase {
public:
    HttpClient(bool bAsyncCall):IHttpBase(bAsyncCall) {}
    ~HttpClient() {}

	virtual void process_http_response(bool bSuccess, const char* response, int status) {
		if(m_bAsyncCall) {
			if(bSuccess) {
				cout << "Async call success, status: " << status << ", response: " << response << endl;
			}
			else {
				cout << "Async call failed, status: " << status << ", response: " << response << endl;
			}
		}
		else {
			if(bSuccess) {
				cout << "Sync call success, status: " << status << ", response: " << response << endl;
			}
			else {
				cout << "Sync call failed, status: " << status << ", response: " << response << endl;
			}
		}
	}
};

int main() {
    string strUrl = "http://192.168.100.7:3901/getdbAddress";
    string strMethod = "GET";
    string strHeaders = "{\"Content-Type\": \"application/json\", \"User-Agent\": \"MyApp\"}";
    string strBody = "Hello 9100 port";
    
#if 1
	printf("开始发送同步HTTP请求...\n");
    HttpClient httpClient(false);
#else
    HttpClient httpClient(true);
#endif

	httpClient.send_http_message(strUrl, strMethod, strHeaders, strBody);

	while(1) {
		cout << "main while ..." << endl;
		usleep(1000000);
	}

    return 0;
}


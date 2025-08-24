#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "send_http_request.h"

typedef struct _stCallbackParam {
    char caller[256];
    int nSuccess; //0 success  1 failed
    char* pResponse;
}stCallbackParam;

void http_callback(const char* response, int status, int error_type, void* user_data)
{
    if(user_data == NULL) { return; }

    stCallbackParam* pCBParam = NULL;
    pCBParam = (stCallbackParam*)user_data;

    switch (error_type) {
        case ERROR_TYPE_NONE:
            if(status != 200) {
                pCBParam->nSuccess = 1;
            }
            else {
                pCBParam->nSuccess = 0;
            }

            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_CONNECTION_FAILED:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, 连接失败: %s\n", (char*)user_data, response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_TIMEOUT:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, 请求超时: %s\n", (char*)user_data, response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_DNS:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, DNS解析错误: %s\n", response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_TLS:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, TLS/SSL错误: %s\n", response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_READ_RESPONSE:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, 读取响应失败: %s\n", response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        case ERROR_TYPE_OTHER:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, 其他错误: %s\n", response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
        default:
            pCBParam->nSuccess = 1;
            //printf("caller: %s, 未知错误类型: %d, 错误信息: %s\n", error_type, response);
            if(response != NULL) {
               pCBParam->pResponse = (char*)malloc(strlen(response) + 1);
               strcpy(pCBParam->pResponse, response);
               pCBParam->pResponse[strlen(response)] = 0x00;
            }
            break;
    }
}

int main() {
    const char* url = "https://192.168.104.100:9100/list";
    const char* method = "GET";
    const char* headers = "{\"Content-Type\": \"application/json\", \"User-Agent\": \"MyApp\"}";
    const char* body = "Hello 9100 port";
     
    stCallbackParam cbparam;
	memset(cbparam.caller, 0, sizeof(cbparam.caller));
    cbparam.nSuccess = 0;
    cbparam.pResponse = NULL;

#if 0
    // 同步调用，通过回调函数返回结果
    printf("开始发送同步HTTP请求...\n");
    strcpy(cbparam.caller, "This is SendHttpRequestSync");
    SendHttpRequestSync(url, method, headers, body, http_callback, (void*)&cbparam);
#else
    // 异步调用，通过回调函数返回结果
    printf("开始发送异步HTTP请求...\n");
    strcpy(cbparam.caller, "This is SendHttpRequestAsync");
    SendHttpRequestAsync(url, method, headers, body, http_callback, (void*)&cbparam);
#endif

	while(cbparam.pResponse == NULL) {
		printf("main while ...\n");
		sleep(1);
	}

    if(cbparam.nSuccess==0 && cbparam.pResponse != NULL) {
        printf("caller: %s, call success, response: %s\n", cbparam.caller , cbparam.pResponse);
        free(cbparam.pResponse);
        cbparam.pResponse = NULL;
    }

    if(cbparam.nSuccess==1 && cbparam.pResponse != NULL) {
        printf("caller: %s, call failed, response: %s\n", cbparam.caller , cbparam.pResponse);
        free(cbparam.pResponse);
        cbparam.pResponse = NULL;
    }

    return 0;
}


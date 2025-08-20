#ifndef __SEND_HTTP_REQUEST_H__
#define __SEND_HTTP_REQUEST_H__

#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif

// 错误类型常量
typedef enum {
    ERROR_TYPE_NONE = 0,            // 无错误
    ERROR_TYPE_CONNECTION_FAILED,   // 连接失败
    ERROR_TYPE_TIMEOUT,             // 超时
    ERROR_TYPE_TLS,                 // TLS错误
    ERROR_TYPE_DNS,                 // DNS解析错误
    ERROR_TYPE_READ_RESPONSE,       // 读取响应失败
    ERROR_TYPE_OTHER                // 其他错误
} ErrorType;

typedef void (*HttpCallback)(const char* response, int status, int error_type, void* user_data);

// 导出函数, 异步发送http请求
void SendHttpRequestAsync(const char* url, const char* method, const char* headers, 
                         const char* body, HttpCallback callback, void* user_data);


// 导出函数, 同步发送http请求
void SendHttpRequestSync(const char* url, const char* method, const char* headers, 
                         const char* body, HttpCallback callback, void* user_data);

#ifdef __cplusplus
}
#endif

#endif // __SEND_HTTP_REQUEST_H__

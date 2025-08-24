#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include "SM4ECBEncryptDecrypt.h"
using namespace std;

int main()
{
    string strKey = "JeF8U9wHFOMfs2Y8";
    string strPlainData = "世界那么大，我想去看看！hello worle , 秘钥是：JeF8U9wHFOMfs2Y8, 我";
    cout << "strKey: " << strKey << endl;
    cout << "strPlainData: " << strPlainData << endl;

    char errMsg[256] = { 0 };
    
    char* ciphertext = SM4ECBEncrypt(strPlainData.data(), strKey.data(), errMsg, sizeof(errMsg));
    if (ciphertext) {
        printf("Encryption successful, ciphertext: %s\n", ciphertext);

        char* decrypted = SM4ECBDecrypt(ciphertext, strKey.data(), errMsg, sizeof(errMsg));
        if (decrypted) {
            printf("Decrypted text: %s\n", decrypted);
            FreeCString(decrypted);
        } else {
            printf("Decryption failed, errorMsg: %s\n", errMsg);
        }
        
        FreeCString(ciphertext);
    } else {
        printf("Encryption failed, errorMsg: %s\n", errMsg);
    }
    
    return 0;
}

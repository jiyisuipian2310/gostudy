#include <iostream>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <string>
#include "SM4CBCEncryptDecrypt.h"
using namespace std;

int main()
{
    string strKey = "JeF8U9wHFOMfs2Y8";
    string strPlainData = "世界那么大，我想去看看！";
    string strIV = "1234567812345671";
    cout << "strKey: " << strKey << endl;
    cout << "strIV: " << strIV << endl;
    cout << "strPlainData: " << strPlainData << endl;

    char errMsg[256] = { 0 };
    
    char* ciphertext = SM4Encrypt(strPlainData.data(), strKey.data(), strIV.data(), errMsg, sizeof(errMsg));
    if (ciphertext) {
        printf("Encryption successful, ciphertext: %s\n", ciphertext);
        
        char* decrypted = SM4Decrypt(ciphertext, strKey.data(), strIV.data(), errMsg, sizeof(errMsg));
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

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

var ( // 定义全局变量
	g_publicKey  *rsa.PublicKey
	g_privateKey *rsa.PrivateKey
)

func generatePrivateKey(bits int, path string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 将PEM格式的私钥写入文件
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(string(privateKeyPem)))
	if err != nil {
		panic(err)
	}

	return nil
}

func loadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解码PEM格式的私钥
	block, _ := pem.Decode(fileBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("Failed to decode PEM block containing RSA private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// 加密函数
func rsa_encrypt(privateKey *rsa.PrivateKey, message string) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte(message))
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// 解密函数
func rsa_decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) (string, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func main() {
	var err error
	private_key_file := "private_key.pem"

	if false {
		// 生成私钥文件
		err := generatePrivateKey(2048, private_key_file)
		if err != nil {
			fmt.Println("Failed to generate private key:", err)
		}
	}

	if _, err := os.Stat(private_key_file); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s 不存在\n", private_key_file)
		} else {
			fmt.Printf("%s 其他错误\n", private_key_file)
		}
		return
	}

	g_privateKey, err = loadPrivateKeyFromFile(private_key_file)
	if err != nil {
		fmt.Printf("loadPrivateKeyFromFile faield: %v\n", err)
		return
	}

	strPlainData := "世界那么大， 我想去看看。"
	encryptData, err := rsa_encrypt(g_privateKey, strPlainData)
	if err != nil {
		fmt.Printf("rsa_encrypt faield: %v\n", err)
		return
	}

	decryptData, err := rsa_decrypt(g_privateKey, encryptData)
	if err != nil {
		fmt.Printf("rsa_decrypt faield: %v\n", err)
		return
	}

	fmt.Printf("原始数据: %s\n", strPlainData)
	fmt.Printf("解密数据: %s\n", decryptData)
}

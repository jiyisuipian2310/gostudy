package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// bits 生成的公私钥对的位数，一般为 1024 或 2048
// privateKey 生成的私钥  publicKey 生成的公钥
func GenRsaKey(bits int) (privateKey, publicKey string, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	derStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvKey := pem.EncodeToMemory(block)
	puKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(puKey)
	if err != nil {
		return "", "", err
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)

	return string(prvKey), string(pubKey), nil
}

// originalData 签名前的原始数据
// privateKey RSA 私钥
func signBase64(originalData, privateKey string) (string, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("解析私钥失败: 无法解码PEM数据")
	}

	if block.Type == "RSA PRIVATE KEY" { // PKCS#1
		fmt.Println("This is PKCS#1 RSA PRIVATE KEY")
		priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("解析私钥失败: %v", err)
		}
		// sha256 加密方式，必须与 下面的 crypto.SHA256 对应
		// 例如使用 sha1 加密，此处应是 sha1.Sum()，对应 crypto.SHA1
		hash := sha256.Sum256([]byte(originalData))
		signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hash[:])
		if err != nil {
			return "", fmt.Errorf("签名失败: %v", err)
		}

		return base64.StdEncoding.EncodeToString(signature), nil
	} else if block.Type == "PRIVATE KEY" { // PKCS#8
		fmt.Println("This is PKCS#8 PRIVATE KEY")

		priKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", fmt.Errorf("解析私钥失败: %v", err)
		}
		// sha256 加密方式，必须与 下面的 crypto.SHA256 对应
		// 例如使用 sha1 加密，此处应是 sha1.Sum()，对应 crypto.SHA1
		hash := sha256.Sum256([]byte(originalData))
		signature, err := rsa.SignPKCS1v15(rand.Reader, priKey.(*rsa.PrivateKey), crypto.SHA256, hash[:])
		if err != nil {
			return "", fmt.Errorf("签名失败: %v", err)
		}

		return base64.StdEncoding.EncodeToString(signature), nil
	} else {
		return "", errors.New("解析私钥失败: 私钥类型错误")
	}
}

// originalData 签名前的原始数据
// signData Base64 格式的签名串
// pubKey 公钥（需与加密时使用的私钥相对应）
// 返回 true 代表验签通过，反之为不通过
func verifySignWithBase64(originalData, signData, pubKey string) (bool, error) {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return false, fmt.Errorf("签名解码失败: %v", err)
	}

	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return false, errors.New("解析公钥失败: 无法解码PEM数据")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("解析公钥失败: %v", err)
	}

	// sha256 加密方式，必须与 下面的 crypto.SHA256 对应
	// 例如使用 sha1 加密，此处应是 sha1.Sum()，对应 crypto.SHA1
	hash := sha256.Sum256([]byte(originalData))
	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash[:], sign)
	if err != nil {
		return false, fmt.Errorf("验签失败: %v", err)
	}

	return true, nil
}

func RsaEncryptBase64(originalData, publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", errors.New("RsaEncryptBase64 公钥解码失败")
	}

	pubKey, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
	if parseErr != nil {
		return "", fmt.Errorf("RsaEncryptBase64 解析公钥失败: %v", parseErr)
	}

	// 获取密钥长度，计算最大加密块大小
	keySize := pubKey.(*rsa.PublicKey).Size()
	maxEncryptSize := keySize - 11

	// 将原始数据按块大小分段加密
	var encryptedData []byte
	for len(originalData) > 0 {
		segment := originalData
		if len(segment) > maxEncryptSize {
			segment = originalData[:maxEncryptSize]
		}

		encryptedSegment, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(segment))
		if err != nil {
			return "", fmt.Errorf("RsaEncryptBase64 加密失败: %v", err)
		}

		encryptedData = append(encryptedData, encryptedSegment...)
		originalData = originalData[len(segment):]
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func RsaDecryptBase64(encryptedData, privateKey string) (string, error) {
	encryptedDecodeBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("RsaDecryptBase64 解码失败: %v", err)
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", errors.New("RsaDecryptBase64 私钥解码失败")
	}

	keySize := 0
	var priKey *rsa.PrivateKey = nil
	var parseErr error

	if block.Type == "RSA PRIVATE KEY" { // PKCS#1
		priKey, parseErr = x509.ParsePKCS1PrivateKey(block.Bytes)
		if parseErr != nil {
			return "", fmt.Errorf("RsaDecryptBase64 解析私钥失败: %v", parseErr)
		}
	} else if block.Type == "PRIVATE KEY" { // PKCS#8
		var key any
		key, parseErr = x509.ParsePKCS8PrivateKey(block.Bytes)
		priKey = key.(*rsa.PrivateKey)
		if parseErr != nil {
			return "", fmt.Errorf("RsaDecryptBase64 解析私钥失败: %v", parseErr)
		}
	} else {
		return "", errors.New("RsaDecryptBase64 私钥类型错误")
	}

	// 获取密钥长度，计算最大解密块大小
	keySize = priKey.Size()

	// 分段解密数据
	var decryptedData []byte
	for len(encryptedDecodeBytes) > 0 {
		segment := encryptedDecodeBytes
		if len(segment) > keySize {
			segment = encryptedDecodeBytes[:keySize]
		}

		decryptedSegment, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, segment)
		if err != nil {
			return "", fmt.Errorf("RsaDecryptBase64 解密失败: %v", err)
		}

		decryptedData = append(decryptedData, decryptedSegment...)
		encryptedDecodeBytes = encryptedDecodeBytes[len(segment):]
	}

	return string(decryptedData), nil
}

func main() {
	privateKey, publicKey, _ := GenRsaKey(2048)

	//签名
	strSignedData, err := signBase64("hello world", privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("SignedData: %s\n", strSignedData)

	//验签
	bResult, err := verifySignWithBase64("hello world", strSignedData, publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("验签结果: %t\n", bResult)

	//公钥加密
	strPlainData := "我想看看这个大千世界"
	strEncryptedData, err := RsaEncryptBase64(strPlainData, publicKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("公钥加密后数据: %s\n", strEncryptedData)

	//私钥解密
	strDecryptedData, err := RsaDecryptBase64(strEncryptedData, privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("私钥解密后数据: %s\n", strDecryptedData)
}

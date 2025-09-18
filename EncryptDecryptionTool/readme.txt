创建加解密命令行工具的步骤：
1. 安装cobra-cli库  ==> go install github.com/spf13/cobra-cli@latest 
2. 使用 Cobra CLI 创建 EncryptDecryptionTool 项目
	mkdir EncryptDecryptionTool
	cd EncryptDecryptionTool
	go mod init EncryptDecryptionTool
	cobra-cli init
	
3. 添加子命令
	cobra-cli add aes_encrypt
	cobra-cli add aes_decrypt
	cobra-cli add sm4_encrypt
	cobra-cli add sm4_decrypt

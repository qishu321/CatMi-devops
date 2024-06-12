package tools

import (
	"CatMi-devops/config"
	"crypto/aes"
	"crypto/cipher"
	rand2 "crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"

	"time"
)

// 在程序启动时设置种子，然后在整个程序运行期间使用这个种子生成随机数。
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RSA密码加密
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// RSA密码解密
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}
func GenerateRandomNumber() int64 {
	// // 设置随机种子，保证每次生成的随机数都不同
	// rand.Seed(time.Now().UnixNano())

	// 生成10位随机整数，范围在0到9999999999之间
	return int64(rand.Intn(10000000000))
}

func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand2.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCipherText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCipherText) < aes.BlockSize {
		return "", fmt.Errorf("加密数据长度错误")
	}

	iv := decodedCipherText[:aes.BlockSize]
	decodedCipherText = decodedCipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCipherText, decodedCipherText)

	return string(decodedCipherText), nil
}

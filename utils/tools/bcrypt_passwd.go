package tools

import (
	"CatMi-devops/config"
	"math/rand"
	"time"
)

// 密码加密
func NewGenPasswd(passwd string) string {
	pass, _ := RSAEncrypt([]byte(passwd), config.Conf.System.RSAPublicBytes)
	return string(pass)
}

// 密码解密
func NewParPasswd(passwd string) string {
	pass, _ := RSADecrypt([]byte(passwd), config.Conf.System.RSAPrivateBytes)
	return string(pass)
}
func GenerateRandomNumber() int64 {
	// 设置随机种子，保证每次生成的随机数都不同
	rand.Seed(time.Now().UnixNano())

	// 生成10位随机整数，范围在0到9999999999之间
	return int64(rand.Intn(10000000000))
}
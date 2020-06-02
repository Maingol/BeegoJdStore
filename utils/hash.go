package utils

import (
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/bcrypt"
)

// 对用户密码进行hash加密、解密，详情参见：https://www.kancloud.cn/golang_programe/golang/1144844

// 根据用户的密码字符串，生成哈希字符串
func HashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return string(hash)
}

//验证密码
//CompareHashAndPassword 将 bcrypt 哈希密码与其纯文本进行比较。 成功时返回 nil，失败时返回错误
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}

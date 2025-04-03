package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MD5Encode 对输入的字符串进行MD5哈希编码。
// 该函数接受一个字符串作为参数，返回该字符串的MD5哈希值的十六进制表示。
// MD5是一种广泛使用的加密算法，尽管它不再被认为是安全的，但它仍然用于生成数据的校验和或数字指纹。
// 参数:
//
//	str - 需要进行MD5编码的字符串。
//
// 返回值:
//
//	返回字符串的MD5哈希值的十六进制表示。
func MD5Encode(str string) string {
	// 创建一个新的MD5哈希对象。
	h := md5.New()

	// 将输入的字符串以字节切片的形式写入哈希对象。
	h.Write([]byte(str))

	// 计算哈希值并返回一个包含哈希结果的字节切片。
	tempStr := h.Sum(nil)

	// 将哈希结果的字节切片转换为十六进制字符串并返回。
	return hex.EncodeToString(tempStr)
}

// MD5EncodeUpper 对输入字符串进行MD5编码，并将结果转换为大写字符串。
// 该函数接受一个字符串作为参数，返回其MD5编码的大写字符串形式。
func MD5EncodeUpper(str string) string {
	return strings.ToUpper(MD5Encode(str))
}

// MakePassword 对给定的密码和盐值进行散列处理，生成安全的密码字符串。
// 该函数接受两个参数：password（原始密码字符串）和salt（用于增加密码强度的盐值字符串）。
// 返回值为经过MD5散列处理后的密码字符串。
func MakePassword(password, salt string) string {
	return MD5Encode(password + salt)
}

// ValidPassword 检查密码是否有效。
// 该函数通过将密码与盐值拼接后进行MD5加密，并将加密结果与给定的加密密码进行比较来验证密码的有效性。
// 参数:
//
//	password - 用户输入的原始密码。
//	salt - 在原始密码上添加的盐值，用于增加密码安全性。
//	encodedPassword - 已加密的密码，用于与计算后的加密密码进行比较。
//
// 返回值:
//
//	如果计算后的加密密码与给定的加密密码相同，则返回true，表示密码验证成功；否则返回false。
func ValidPassword(password, salt, encodedPassword string) bool {
	// 比较给定的加密密码与使用盐值对原始密码进行MD5加密后的结果是否一致
	return encodedPassword == MD5Encode(password+salt)
}

// MakeToken 生成一个随机的令牌字符串。
func MakeToken() string {
	return MD5Encode(fmt.Sprintf("%d", time.Now().Unix()))
}

// MakeSalt 生成一个随机的盐值字符串。
func MakeSalt() string {
	return fmt.Sprintf("%06d", rand.Int31()%1000000)
}

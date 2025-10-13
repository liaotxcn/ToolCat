package pkg

import (
	"math/rand"
	"time"
)

var (// 随机字符串的字符集
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	random     = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomString 生成指定长度的随机字符串
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(b)
}

// StrSliceContains 检查字符串切片是否包含指定字符串
func StrSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GenerateUniqueID 生成唯一ID（基于时间戳和随机数）
func GenerateUniqueID() string {
	return time.Now().Format("20060102150405") + "-" + RandomString(8)
}

// StringInSlice 检查字符串是否在切片中（别名函数，为了兼容）
func StringInSlice(str string, list []string) bool {
	return StrSliceContains(list, str)
}
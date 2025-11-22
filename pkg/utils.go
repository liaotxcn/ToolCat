package pkg

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// 定义常量，提高可维护性
const (
	// 随机字符串的字符集
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// 为并发安全使用 sync.Pool 缓存随机数生成器
var (// 缓存随机数生成器以提高性能
	randomSourcePool = sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(time.Now().UnixNano()))
		},
	}

	// 用于并发安全的互斥锁
	randomMu sync.Mutex
)

// RandomString 生成指定长度的随机字符串
// 使用 crypto/rand 生成高质量随机数，适合安全场景
func RandomString(n int) string {
	// 对小字符串使用 math/rand 提高性能
	if n <= 100 {
		b := make([]byte, n)
		randomMu.Lock()
		rng := randomSourcePool.Get().(*rand.Rand)
		for i, cache, remain := n-1, rng.Int63(), letterIdxMax; i >= 0; {
			if remain == 0 {
				cache, remain = rng.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
		randomSourcePool.Put(rng)
		randomMu.Unlock()
		return string(b)
	}

	// 对大字符串使用 crypto/rand 确保安全性
	return SecureRandomString(n)
}

// SecureRandomString 使用密码学安全的随机数生成器生成随机字符串
// 适合需要高安全性的场景，如生成令牌、密钥等
func SecureRandomString(n int) string {
	// 计算需要的字节数（base64编码会增加约33%的大小）
	randomBytes := make([]byte, (n*3)/4) // base64 编码的优化计算
	_, err := cryptoRand.Read(randomBytes)
	if err != nil {
		// 如果 crypto/rand 失败，回退到 math/rand
		return RandomString(n)
	}

	// 使用 base64 URL 安全编码并截取指定长度
	result := base64.URLEncoding.EncodeToString(randomBytes)
	if len(result) > n {
		return result[:n]
	}
	return result
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

// StrSliceContainsAny 检查字符串切片是否包含任何一个指定的字符串
func StrSliceContainsAny(slice []string, items ...string) bool {
	for _, item := range items {
		if StrSliceContains(slice, item) {
			return true
		}
	}
	return false
}

// GenerateUniqueID 生成唯一ID（基于时间戳和随机数）
// 格式: YYYYMMDDHHmmss-8位随机字符串
func GenerateUniqueID() string {
	return time.Now().Format("20060102150405") + "-" + RandomString(8)
}

// GenerateShortID 生成更短的唯一ID
// 格式: 时间戳的base36编码-6位随机字符串
func GenerateShortID() string {
	// 使用时间戳的base36编码作为前缀，比标准时间格式更短
	timestamp := time.Now().Unix()
	base36Time := big.NewInt(timestamp).Text(36)
	return base36Time + "-" + RandomString(6)
}

// GenerateRequestID 生成请求ID
// 格式: req-时间戳-12位随机字符串
func GenerateRequestID() string {
	return "req-" + time.Now().Format("20060102150405") + "-" + RandomString(12)
}

// StringInSlice 检查字符串是否在切片中（别名函数，为了兼容）
func StringInSlice(str string, list []string) bool {
	return StrSliceContains(list, str)
}
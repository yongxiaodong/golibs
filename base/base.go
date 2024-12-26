package base

import (
	"math/rand"
	"os"
	"time"
	"unicode/utf8"
)

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Obfuscate 用户名模糊化, 格式《会***猪》《1***9》
func Obfuscate(input string) string {
	runes := []rune(input)
	count := utf8.RuneCountInString(input)
	if count == 0 {
		return "***"
	}
	// 11位数字则认为是手机号，返回前3位和后四位
	if count == 11 && isAllDigits(input) {
		return string(runes[:3]) + "***" + string(runes[7:])
	}
	if count <= 2 {
		return string(runes[0]) + "***"
	}

	return string(runes[0]) + "***" + string(runes[len(runes)-1])
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// 校验目录，不存在则创建
func CreateDirIfNotExists(path string) error {
	if !dirExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// 生成特定长度的随机字符
func RandString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, n)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

package tinyurl

import (
	"crypto/md5"
	"fmt"
)

var runes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(string) string {
	return ""
}

func GeneragteMD5Value(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

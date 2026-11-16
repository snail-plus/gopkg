package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// HmacSHA256 使用 key 计算 message 的 HMAC-SHA256，返回十六进制字符串。
func HmacSHA256(key, message string) (string, error) {
	mac := hmac.New(sha256.New, []byte(key))
	if _, err := io.WriteString(mac, message); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func MustHmacSHA256(key, message string) string {
	return mustHashString(HmacSHA256(key, message))
}

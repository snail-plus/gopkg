package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// RandomBytes 返回 n 个密码学安全的随机字节。
func RandomBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, fmt.Errorf("crypto: invalid random byte length %d", n)
	}
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// RandomToken 生成 nBytes 个随机字节并编码为 URL 安全的 Base64 字符串。
func RandomToken(nBytes int) (string, error) {
	if nBytes <= 0 {
		return "", fmt.Errorf("crypto: token byte length must be positive, got %d", nBytes)
	}
	b, err := RandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func MustRandomBytes(n int) []byte {
	b, err := RandomBytes(n)
	if err != nil {
		panic(err)
	}
	return b
}

func MustRandomToken(nBytes int) string {
	token, err := RandomToken(nBytes)
	if err != nil {
		panic(err)
	}
	return token
}

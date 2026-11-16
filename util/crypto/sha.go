package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
)

// Sha224 计算字符串的 SHA-224 摘要，返回十六进制字符串。
func Sha224(input string) (string, error) {
	return hashString(sha256.New224(), input)
}

// Sha256 计算字符串的 SHA-256 摘要，返回十六进制字符串。
func Sha256(input string) (string, error) {
	return hashString(sha256.New(), input)
}

// Sha384 计算字符串的 SHA-384 摘要，返回十六进制字符串。
func Sha384(input string) (string, error) {
	return hashString(sha512.New384(), input)
}

// Sha512 计算字符串的 SHA-512 摘要，返回十六进制字符串。
func Sha512(input string) (string, error) {
	return hashString(sha512.New(), input)
}

func hashString(h hash.Hash, input string) (string, error) {
	if _, err := io.WriteString(h, input); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func mustHashString(res string, err error) string {
	if err != nil {
		panic(err)
	}
	return res
}

func MustSha224(input string) string { return mustHashString(Sha224(input)) }

func MustSha256(input string) string { return mustHashString(Sha256(input)) }

func MustSha384(input string) string { return mustHashString(Sha384(input)) }

func MustSha512(input string) string { return mustHashString(Sha512(input)) }

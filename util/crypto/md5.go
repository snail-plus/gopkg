package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"strings"
)

// Md5 计算字符串的 MD5 摘要，返回十六进制字符串。
// 注意：MD5 已不具备抗碰撞性，不应用于安全相关场景，仅用于兼容性校验。
func Md5(input string) (string, error) {
	return Md5Reader(strings.NewReader(input))
}

// Md5Reader 计算 io.Reader 中数据的 MD5 摘要，返回十六进制字符串。
func Md5Reader(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Md5File 计算文件的 MD5 摘要，返回十六进制字符串。
func Md5File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return Md5Reader(f)
}

func MustMd5(input string) string {
	res, err := Md5(input)
	if err != nil {
		panic(err)
	}
	return res
}

func MustMd5Reader(r io.Reader) string {
	res, err := Md5Reader(r)
	if err != nil {
		panic(err)
	}
	return res
}

func MustMd5File(path string) string {
	res, err := Md5File(path)
	if err != nil {
		panic(err)
	}
	return res
}

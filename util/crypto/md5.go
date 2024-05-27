package crypto

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(input string) (string, error) {

	// 创建一个新的 hash.Hash 接口来写入数据
	md5Hash := md5.New()
	// 将字符串转换为字节并写入 hasher
	_, err := io.WriteString(md5Hash, input)
	if err != nil {
		return "", err
	}

	// 计算 MD5 哈希值并获取字节切片
	sum := md5Hash.Sum(nil)

	// 将字节切片转换为十六进制字符串
	md5String := hex.EncodeToString(sum)
	return md5String, nil
}

func MustMd5(input string) string {
	res, err := Md5(input)
	if err != nil {
		panic(err)
	}
	return res
}

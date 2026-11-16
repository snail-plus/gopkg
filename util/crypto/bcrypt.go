package crypto

import "golang.org/x/crypto/bcrypt"

// BcryptHash 使用默认 cost 对密码做 bcrypt 哈希。
// 注意 bcrypt 最多处理 72 字节的密码，超长会返回错误。
func BcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BcryptCompare 校验明文密码与 bcrypt 哈希是否匹配，匹配时返回 nil。
func BcryptCompare(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

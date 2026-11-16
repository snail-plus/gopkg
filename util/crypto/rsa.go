package crypto

import (
	stdcrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// rsaMinBits 是生成密钥时允许的最小位数，低于该值视为不安全。
const rsaMinBits = 1024

var (
	ErrInvalidPEM       = errors.New("crypto: invalid PEM data")
	ErrNotRSAPrivateKey = errors.New("crypto: not an RSA private key")
	ErrNotRSAPublicKey  = errors.New("crypto: not an RSA public key")
)

// GenerateRSAKeyPair 生成 RSA 密钥对。bits 不得小于 1024，生产环境建议 2048 或以上。
func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	if bits < rsaMinBits {
		return nil, nil, fmt.Errorf("crypto: rsa key size must be at least %d bits, got %d", rsaMinBits, bits)
	}
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// RSAEncrypt 使用 OAEP（SHA-256）加密 plaintext。
// 明文长度上限为密钥字节数减去 2*sha256.Size+2，超长会返回错误。
func RSAEncrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	if pub == nil {
		return nil, errors.New("crypto: nil rsa public key")
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, nil)
}

// RSADecrypt 使用 OAEP（SHA-256）解密由 RSAEncrypt 生成的密文。
func RSADecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("crypto: nil rsa private key")
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
}

// RSASign 对 data 的 SHA-256 摘要做 PKCS#1 v1.5 签名。
func RSASign(priv *rsa.PrivateKey, data []byte) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("crypto: nil rsa private key")
	}
	digest := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, priv, stdcrypto.SHA256, digest[:])
}

// RSAVerify 校验 data 的 PKCS#1 v1.5 签名，签名不匹配时返回错误。
func RSAVerify(pub *rsa.PublicKey, data, signature []byte) error {
	if pub == nil {
		return errors.New("crypto: nil rsa public key")
	}
	digest := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(pub, stdcrypto.SHA256, digest[:], signature)
}

// EncodePrivateKeyPEM 将 RSA 私钥编码为 PKCS#8 PEM。
func EncodePrivateKeyPEM(priv *rsa.PrivateKey) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("crypto: nil rsa private key")
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// EncodePublicKeyPEM 将 RSA 公钥编码为 PKIX PEM。
func EncodePublicKeyPEM(pub *rsa.PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, errors.New("crypto: nil rsa public key")
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

// parseRSAPrivateKey 解析 PEM 编码的 RSA 私钥，兼容 PKCS#8 与 PKCS#1 格式。
func parseRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidPEM
	}
	switch block.Type {
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, ErrNotRSAPrivateKey
		}
		return rsaKey, nil
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, ErrInvalidPEM
	}
}

// parseRSAPublicKey 解析 PEM 编码的 RSA 公钥，兼容 PKIX、PKCS#1 与 X.509 证书。
func parseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidPEM
	}
	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, ErrNotRSAPublicKey
		}
		return rsaKey, nil
	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)
	case "CERTIFICATE":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, ErrNotRSAPublicKey
		}
		return rsaKey, nil
	default:
		return nil, ErrInvalidPEM
	}
}

// DecodePrivateKeyPEM 解析 PEM 格式的 RSA 私钥，兼容 PKCS#8 与 PKCS#1。
func DecodePrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	return parseRSAPrivateKey(data)
}

// DecodePublicKeyPEM 解析 PEM 格式的 RSA 公钥，兼容 PKIX、PKCS#1 与 X.509 证书。
func DecodePublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	return parseRSAPublicKey(data)
}

// ParseRSAPrivateKey 从 PEM 字符串解析 RSA 私钥。
func ParseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	priv, err := parseRSAPrivateKey([]byte(strings.TrimSpace(pemStr)))
	if err != nil {
		return nil, fmt.Errorf("crypto: parse rsa private key: %w", err)
	}
	return priv, nil
}

// ParseRSAPublicKey 从 PEM 字符串解析 RSA 公钥。
func ParseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	pub, err := parseRSAPublicKey([]byte(strings.TrimSpace(pemStr)))
	if err != nil {
		return nil, fmt.Errorf("crypto: parse rsa public key: %w", err)
	}
	return pub, nil
}

// RSAEncryptWithPEM 使用 PEM 字符串形式的公钥加密 plaintext。
func RSAEncryptWithPEM(pubPEM string, plaintext []byte) ([]byte, error) {
	pub, err := ParseRSAPublicKey(pubPEM)
	if err != nil {
		return nil, err
	}
	return RSAEncrypt(pub, plaintext)
}

// RSADecryptWithPEM 使用 PEM 字符串形式的私钥解密密文。
func RSADecryptWithPEM(privPEM string, ciphertext []byte) ([]byte, error) {
	priv, err := ParseRSAPrivateKey(privPEM)
	if err != nil {
		return nil, err
	}
	return RSADecrypt(priv, ciphertext)
}

// RSASignWithPEM 使用 PEM 字符串形式的私钥对 data 签名。
func RSASignWithPEM(privPEM string, data []byte) ([]byte, error) {
	priv, err := ParseRSAPrivateKey(privPEM)
	if err != nil {
		return nil, err
	}
	return RSASign(priv, data)
}

// RSAVerifyWithPEM 使用 PEM 字符串形式的公钥校验签名。
func RSAVerifyWithPEM(pubPEM string, data, signature []byte) error {
	pub, err := ParseRSAPublicKey(pubPEM)
	if err != nil {
		return err
	}
	return RSAVerify(pub, data, signature)
}

// RSAEncryptText 使用 PEM 字符串形式的公钥加密文本，返回 Base64 编码的密文。
func RSAEncryptText(pubPEM, text string) (string, error) {
	ciphertext, err := RSAEncryptWithPEM(pubPEM, []byte(text))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSADecryptText 解密 Base64 编码的密文并返回明文文本。
func RSADecryptText(privPEM, encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("crypto: decode rsa ciphertext: %w", err)
	}
	plaintext, err := RSADecryptWithPEM(privPEM, ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// RSASignText 使用 PEM 字符串形式的私钥对文本签名，返回 Base64 编码的签名。
func RSASignText(privPEM, text string) (string, error) {
	signature, err := RSASignWithPEM(privPEM, []byte(text))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// RSAVerifyText 校验 Base64 编码的签名是否与文本匹配。
func RSAVerifyText(pubPEM, text, signature string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("crypto: decode rsa signature: %w", err)
	}
	return RSAVerifyWithPEM(pubPEM, []byte(text), sig)
}

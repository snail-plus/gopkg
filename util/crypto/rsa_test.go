package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"testing"
)

func TestRSAEncryptDecrypt(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("hello rsa")
	ciphertext, err := RSAEncrypt(pub, msg)
	if err != nil {
		t.Fatal(err)
	}
	got, err := RSADecrypt(priv, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(msg) {
		t.Fatalf("got %q, want %q", got, msg)
	}

	otherPriv, _, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := RSADecrypt(otherPriv, ciphertext); err == nil {
		t.Fatal("decrypt with wrong private key expected error")
	}
}

func TestRSASignVerify(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("message to sign")
	signature, err := RSASign(priv, data)
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerify(pub, data, signature); err != nil {
		t.Fatalf("verify valid signature: %v", err)
	}
	if err := RSAVerify(pub, []byte("tampered"), signature); err == nil {
		t.Fatal("verify tampered data expected error")
	}
}

func TestRSAPEMRoundTrip(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}

	privPEM, err := EncodePrivateKeyPEM(priv)
	if err != nil {
		t.Fatal(err)
	}
	gotPriv, err := DecodePrivateKeyPEM(privPEM)
	if err != nil {
		t.Fatal(err)
	}
	if gotPriv.D.Cmp(priv.D) != 0 || gotPriv.N.Cmp(priv.N) != 0 {
		t.Fatal("private key PEM round trip mismatch")
	}

	pubPEM, err := EncodePublicKeyPEM(pub)
	if err != nil {
		t.Fatal(err)
	}
	gotPub, err := DecodePublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatal(err)
	}
	if gotPub.E != pub.E || gotPub.N.Cmp(pub.N) != 0 {
		t.Fatal("public key PEM round trip mismatch")
	}
}

func TestRSAInvalidPEM(t *testing.T) {
	if _, err := DecodePrivateKeyPEM([]byte("not pem")); !errors.Is(err, ErrInvalidPEM) {
		t.Fatalf("got %v, want ErrInvalidPEM", err)
	}
	if _, err := DecodePublicKeyPEM([]byte("not pem")); !errors.Is(err, ErrInvalidPEM) {
		t.Fatalf("got %v, want ErrInvalidPEM", err)
	}
}

func TestGenerateRSAKeyPairTooSmall(t *testing.T) {
	if _, _, err := GenerateRSAKeyPair(512); err == nil {
		t.Fatal("rsa key size below 1024 expected error")
	}
}

func TestRSAConvenienceWithPEMString(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	privPEM, err := EncodePrivateKeyPEM(priv)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM, err := EncodePublicKeyPEM(pub)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("pem key message")
	ciphertext, err := RSAEncryptWithPEM(string(pubPEM), msg)
	if err != nil {
		t.Fatal(err)
	}
	got, err := RSADecryptWithPEM(string(privPEM), ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(msg) {
		t.Fatalf("got %q, want %q", got, msg)
	}

	signature, err := RSASignWithPEM(string(privPEM), msg)
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerifyWithPEM(string(pubPEM), msg, signature); err != nil {
		t.Fatalf("verify valid signature: %v", err)
	}
	if err := RSAVerifyWithPEM(string(pubPEM), []byte("tampered"), signature); err == nil {
		t.Fatal("verify tampered data expected error")
	}
}

func TestRSAConvenienceWithText(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	privPEM, err := EncodePrivateKeyPEM(priv)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM, err := EncodePublicKeyPEM(pub)
	if err != nil {
		t.Fatal(err)
	}
	pubStr := string(pubPEM) + "\n\n"

	text := "hello 文本"
	encoded, err := RSAEncryptText(pubStr, text)
	if err != nil {
		t.Fatal(err)
	}
	got, err := RSADecryptText(string(privPEM), encoded)
	if err != nil {
		t.Fatal(err)
	}
	if got != text {
		t.Fatalf("got %q, want %q", got, text)
	}

	signature, err := RSASignText(string(privPEM), text)
	if err != nil {
		t.Fatal(err)
	}
	if err := RSAVerifyText(pubStr, text, signature); err != nil {
		t.Fatalf("verify valid signature: %v", err)
	}
	if err := RSAVerifyText(pubStr, "changed", signature); err == nil {
		t.Fatal("verify changed text expected error")
	}

	if _, err := RSADecryptText(string(privPEM), "not base64!!"); err == nil {
		t.Fatal("invalid base64 ciphertext expected error")
	}
	if err := RSAVerifyText(pubStr, text, "not base64!!"); err == nil {
		t.Fatal("invalid base64 signature expected error")
	}
}

func TestRSAParsePKCS1PEM(t *testing.T) {
	priv, pub, err := GenerateRSAKeyPair(1024)
	if err != nil {
		t.Fatal(err)
	}
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}
	pubBlock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(pub)}

	parsedPriv, err := ParseRSAPrivateKey(string(pem.EncodeToMemory(privBlock)))
	if err != nil {
		t.Fatal(err)
	}
	parsedPub, err := ParseRSAPublicKey(string(pem.EncodeToMemory(pubBlock)))
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, err := RSAEncrypt(parsedPub, []byte("pkcs1"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := RSADecrypt(parsedPriv, ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "pkcs1" {
		t.Fatalf("got %q", got)
	}
}

func TestRSAInvalidPEMString(t *testing.T) {
	if _, err := ParseRSAPrivateKey("bad"); !errors.Is(err, ErrInvalidPEM) {
		t.Fatalf("got %v, want ErrInvalidPEM", err)
	}
	if _, err := ParseRSAPublicKey("bad"); !errors.Is(err, ErrInvalidPEM) {
		t.Fatalf("got %v, want ErrInvalidPEM", err)
	}
}

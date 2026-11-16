package crypto

import (
	"errors"
	"testing"
)

func TestAESRoundTrip(t *testing.T) {
	plaintext := []byte("hello AES-GCM 中文")
	for _, keyLen := range []int{16, 24, 32} {
		key := make([]byte, keyLen)
		ciphertext, err := AESEncrypt(key, plaintext)
		if err != nil {
			t.Fatalf("AESEncrypt with %d-byte key: %v", keyLen, err)
		}
		got, err := AESDecrypt(key, ciphertext)
		if err != nil {
			t.Fatalf("AESDecrypt with %d-byte key: %v", keyLen, err)
		}
		if string(got) != string(plaintext) {
			t.Fatalf("round trip mismatch: got %q", got)
		}
	}
}

func TestAESInvalidKey(t *testing.T) {
	if _, err := AESEncrypt([]byte("short"), []byte("data")); err == nil {
		t.Fatal("AESEncrypt with invalid key size expected error")
	}
}

func TestAESTamperDetection(t *testing.T) {
	key := make([]byte, 16)
	ciphertext, err := AESEncrypt(key, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	ciphertext[len(ciphertext)-1] ^= 0x01
	if _, err := AESDecrypt(key, ciphertext); err == nil {
		t.Fatal("tampered ciphertext should fail authentication")
	}
}

func TestAESDecryptShortCiphertext(t *testing.T) {
	_, err := AESDecrypt(make([]byte, 16), []byte{1, 2, 3})
	if !errors.Is(err, ErrInvalidCiphertext) {
		t.Fatalf("got %v, want ErrInvalidCiphertext", err)
	}
}

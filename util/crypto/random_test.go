package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	a, err := RandomBytes(32)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 32 {
		t.Fatalf("len = %d, want 32", len(a))
	}
	b := MustRandomBytes(32)
	if bytes.Equal(a, b) {
		t.Fatal("two random byte slices should differ")
	}
	if _, err := RandomBytes(-1); err == nil {
		t.Fatal("negative length expected error")
	}
}

func TestRandomToken(t *testing.T) {
	token, err := RandomToken(32)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		t.Fatalf("token is not URL-safe base64: %v", err)
	}
	if len(decoded) != 32 {
		t.Fatalf("decoded len = %d, want 32", len(decoded))
	}
	if MustRandomToken(32) == token {
		t.Fatal("two random tokens should differ")
	}
	if _, err := RandomToken(0); err == nil {
		t.Fatal("zero byte length expected error")
	}
}

package crypto

import "testing"

func TestHmacSHA256(t *testing.T) {
	want := "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8"
	got := MustHmacSHA256("key", "The quick brown fox jumps over the lazy dog")
	if got != want {
		t.Fatalf("HmacSHA256 = %q, want %q", got, want)
	}

	if _, err := HmacSHA256("key", "message"); err != nil {
		t.Fatalf("HmacSHA256 error: %v", err)
	}
	if MustHmacSHA256("other", "message") == MustHmacSHA256("key", "message") {
		t.Fatal("different keys produced the same HMAC")
	}
}

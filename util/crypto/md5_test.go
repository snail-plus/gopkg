package crypto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMd5(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"abc", "abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"hello", "hello", "5d41402abc4b2a76b9719d911017c592"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustMd5(tt.input); got != tt.want {
				t.Fatalf("Md5(%q) = %q, want %q", tt.input, got, tt.want)
			}
			if _, err := Md5(tt.input); err != nil {
				t.Fatalf("Md5(%q) error: %v", tt.input, err)
			}
		})
	}
}

func TestMd5Reader(t *testing.T) {
	want := MustMd5("abc")
	if got := MustMd5Reader(strings.NewReader("abc")); got != want {
		t.Fatalf("Md5Reader = %q, want %q", got, want)
	}
}

func TestMd5File(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(path, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}

	want := MustMd5("abc")
	if got, err := Md5File(path); err != nil || got != want {
		t.Fatalf("Md5File = %q, %v; want %q", got, err, want)
	}
	if got := MustMd5File(path); got != want {
		t.Fatalf("MustMd5File = %q, want %q", got, want)
	}

	if _, err := Md5File(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatal("Md5File on missing file expected error")
	}
}

package crypto

import "testing"

func TestSha(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"sha224 abc", MustSha224("abc"), "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7"},
		{"sha256 abc", MustSha256("abc"), "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{"sha384 abc", MustSha384("abc"), "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7"},
		{"sha512 abc", MustSha512("abc"), "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f"},
		{"sha256 empty", MustSha256(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}

	if _, err := Sha256("abc"); err != nil {
		t.Fatalf("Sha256 error: %v", err)
	}
}

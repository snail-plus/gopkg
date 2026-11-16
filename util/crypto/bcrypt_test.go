package crypto

import "testing"

func TestBcryptHashAndCompare(t *testing.T) {
	hash, err := BcryptHash("s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if err := BcryptCompare("s3cret", hash); err != nil {
		t.Fatalf("correct password should match: %v", err)
	}
	if err := BcryptCompare("wrong", hash); err == nil {
		t.Fatal("wrong password should not match")
	}

	hash2, err := BcryptHash("s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if hash == hash2 {
		t.Fatal("bcrypt hashes should differ because of random salt")
	}
}

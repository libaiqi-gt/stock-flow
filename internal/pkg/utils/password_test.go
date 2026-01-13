package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "123456"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed: %v", err)
	}
	if hash == password {
		t.Error("Hash should not equal password")
	}
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash failed")
	}
	if CheckPasswordHash("wrong", hash) {
		t.Error("CheckPasswordHash should fail for wrong password")
	}
}

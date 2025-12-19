package cipher

import (
	"fmt"
	"testing"
)

func TestEncryptDecrypt_SameKey(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef" 
	tests := []struct {
		plaintext string
	}{
		{"hello world"},
		{"secret data"},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("plaintext=%s", tt.plaintext)
		t.Run(name, func(t *testing.T) {
			ciphertext, err := Encrypt(key, tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			got, err := Decrypt(key, ciphertext)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			if got != tt.plaintext {
				t.Errorf("Decrypt() = %q; want %q", got, tt.plaintext)
			}
		})
	}
}

func TestDecrypt_DifferentKey_Error(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	otherKey := "abcdef0123456789abcdef0123456789"

	plaintext := "sensitive info"

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	if _, err := Decrypt(otherKey, ciphertext); err == nil {
		t.Errorf("Decrypt() with different key = nil error; want error")
	}
}

func TestEncrypt_InvalidKey_Error(t *testing.T) {
	invalidKey := "short-key" 
	plaintext := "secret data"

	if _, err := Encrypt(invalidKey, plaintext); err == nil {
		t.Errorf("Encrypt() with invalid key = nil error; want error")
	}
}


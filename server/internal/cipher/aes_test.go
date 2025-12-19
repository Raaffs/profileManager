package cipher

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"testing"
)



func GenerateAES256KeyBase64() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func TestEncryptDecrypt_SameKey(t *testing.T) {
	key,err := GenerateAES256KeyBase64()
	if err != nil {
		t.Fatalf("GenerateAES256KeyBase64() error = %v", err)
	}
	fmt.Println(key)
	tests := []struct {
		plaintext string
	}{
		{"hello world"},
		{"secretdata"},
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
	key,err := GenerateAES256KeyBase64()
	if err != nil {
		t.Fatalf("GenerateAES256KeyBase64() error = %v", err)
	}
	otherKey,err:=GenerateAES256KeyBase64()
	if err != nil {
		t.Fatalf("GenerateAES256KeyBase64() error = %v", err)
	}
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
	invalidKey := "fefefwefwf" //not 32 bytes base64 encoded
	plaintext := "secret data"

	if _, err := Encrypt(invalidKey, plaintext); err == nil {
		t.Errorf("Encrypt() with invalid key = nil error; want error")
	}
}


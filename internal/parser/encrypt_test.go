package parser

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	plain := "super-secret-value"
	enc, err := Encrypt(plain, "passphrase123")
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}
	if enc == plain {
		t.Fatal("expected ciphertext to differ from plaintext")
	}
	dec, err := Decrypt(enc, "passphrase123")
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}
	if dec != plain {
		t.Fatalf("expected %q, got %q", plain, dec)
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := Encrypt("value", "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	enc, _ := Encrypt("value", "correct")
	_, err := Decrypt(enc, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("!!!notbase64!!!", "pass")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestEncryptEntries_OnlySensitiveKeys(t *testing.T) {
	entries := []EnvEntry{
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	out, err := EncryptEntries(entries, "pass", DefaultEncryptOptions())
	if err != nil {
		t.Fatalf("EncryptEntries error: %v", err)
	}
	if !strings.HasPrefix(out[0].Value, "enc:") {
		t.Errorf("expected sensitive key to be encrypted, got %q", out[0].Value)
	}
	if out[1].Value != "myapp" {
		t.Errorf("expected non-sensitive key unchanged, got %q", out[1].Value)
	}
}

func TestDecryptEntries_RestoresValues(t *testing.T) {
	entries := []EnvEntry{
		{Key: "API_SECRET", Value: "topsecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	encrypted, _ := EncryptEntries(entries, "pass", DefaultEncryptOptions())
	decrypted, err := DecryptEntries(encrypted, "pass")
	if err != nil {
		t.Fatalf("DecryptEntries error: %v", err)
	}
	if decrypted[0].Value != "topsecret" {
		t.Errorf("expected decrypted value %q, got %q", "topsecret", decrypted[0].Value)
	}
	if decrypted[1].Value != "myapp" {
		t.Errorf("expected unchanged value %q, got %q", "myapp", decrypted[1].Value)
	}
}

func TestEncryptEntries_SkipsEmptyValues(t *testing.T) {
	entries := []EnvEntry{
		{Key: "API_SECRET", Value: ""},
	}
	out, err := EncryptEntries(entries, "pass", DefaultEncryptOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "" {
		t.Errorf("expected empty value to remain empty, got %q", out[0].Value)
	}
}

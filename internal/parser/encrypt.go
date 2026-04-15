package parser

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptOptions controls encryption behaviour.
type EncryptOptions struct {
	KeyEnvVar string // name of the env var holding the passphrase
}

// DefaultEncryptOptions returns sensible defaults.
func DefaultEncryptOptions() EncryptOptions {
	return EncryptOptions{KeyEnvVar: "ENVOY_SECRET_KEY"}
}

// deriveKey stretches a passphrase into a 32-byte AES-256 key.
func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

// Encrypt encrypts plaintext using AES-GCM and returns a base64-encoded ciphertext.
func Encrypt(plaintext, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("encrypt: passphrase must not be empty")
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext.
func Decrypt(encoded, passphrase string) (string, error) {
	if passphrase == "" {
		return "", errors.New("decrypt: passphrase must not be empty")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", errors.New("decrypt: invalid base64 input")
	}
	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("decrypt: ciphertext too short")
	}
	plaintext, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", errors.New("decrypt: authentication failed — wrong passphrase or corrupted data")
	}
	return string(plaintext), nil
}

// EncryptEntries returns a new slice of EnvEntry where sensitive values are encrypted.
func EncryptEntries(entries []EnvEntry, passphrase string, opts EncryptOptions) ([]EnvEntry, error) {
	out := make([]EnvEntry, len(entries))
	for i, e := range entries {
		out[i] = e
		if IsSensitive(e.Key) && e.Value != "" {
			enc, err := Encrypt(e.Value, passphrase)
			if err != nil {
				return nil, err
			}
			out[i].Value = "enc:" + enc
		}
	}
	return out, nil
}

// DecryptEntries returns a new slice of EnvEntry with encrypted values decrypted.
func DecryptEntries(entries []EnvEntry, passphrase string) ([]EnvEntry, error) {
	out := make([]EnvEntry, len(entries))
	for i, e := range entries {
		out[i] = e
		if len(e.Value) > 4 && e.Value[:4] == "enc:" {
			plain, err := Decrypt(e.Value[4:], passphrase)
			if err != nil {
				return nil, err
			}
			out[i].Value = plain
		}
	}
	return out, nil
}

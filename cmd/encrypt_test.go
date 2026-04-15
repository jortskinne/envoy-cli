package cmd

import (
	"os"
	"strings"
	"testing"

	"envoy-cli/internal/parser"
)

func writeEncryptTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestEncryptCmd_RoundTrip(t *testing.T) {
	src := writeEncryptTempEnv(t, "API_SECRET=topsecret\nAPP_NAME=myapp\n")
	out := src + ".enc"
	t.Setenv("ENVOY_SECRET_KEY", "testpass")

	encryptOutput = out
	encryptPassphrase = ""
	err := runEncrypt(encryptCmd, []string{src})
	if err != nil {
		t.Fatalf("runEncrypt error: %v", err)
	}

	encrypted, err := parser.ParseFile(out)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	var secretVal string
	for _, e := range encrypted {
		if e.Key == "API_SECRET" {
			secretVal = e.Value
		}
	}
	if !strings.HasPrefix(secretVal, "enc:") {
		t.Fatalf("expected encrypted value, got %q", secretVal)
	}

	decOut := src + ".dec"
	encryptOutput = decOut
	err = runDecrypt(decryptCmd, []string{out})
	if err != nil {
		t.Fatalf("runDecrypt error: %v", err)
	}

	decrypted, err := parser.ParseFile(decOut)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}
	for _, e := range decrypted {
		if e.Key == "API_SECRET" && e.Value != "topsecret" {
			t.Errorf("expected decrypted value %q, got %q", "topsecret", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected unchanged value %q, got %q", "myapp", e.Value)
		}
	}
}

func TestEncryptCmd_MissingPassphrase(t *testing.T) {
	src := writeEncryptTempEnv(t, "API_SECRET=value\n")
	encryptPassphrase = ""
	_ = os.Unsetenv("ENVOY_SECRET_KEY")
	err := runEncrypt(encryptCmd, []string{src})
	if err == nil {
		t.Fatal("expected error when passphrase is missing")
	}
}

func TestEncryptCmd_PassphraseFlag(t *testing.T) {
	src := writeEncryptTempEnv(t, "DB_PASSWORD=secret\n")
	out := src + ".out"
	encryptOutput = out
	encryptPassphrase = "flagpass"
	_ = os.Unsetenv("ENVOY_SECRET_KEY")
	err := runEncrypt(encryptCmd, []string{src})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, _ := parser.ParseFile(out)
	for _, e := range entries {
		if e.Key == "DB_PASSWORD" && !strings.HasPrefix(e.Value, "enc:") {
			t.Errorf("expected encrypted value, got %q", e.Value)
		}
	}
}

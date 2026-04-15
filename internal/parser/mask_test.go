package parser

import "testing"

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key       string
		expected  bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"AUTH_TOKEN", true},
		{"APP_ENV", false},
		{"PORT", false},
		{"database_secret", true},
		{"PRIVATE_KEY_PATH", true},
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := IsSensitive(tc.key, DefaultSecretPatterns)
			if got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestMaskValue_FullMask(t *testing.T) {
	opts := DefaultMaskOptions()
	result := MaskValue("supersecret", opts)
	if result != "****" {
		t.Errorf("expected '****', got %q", result)
	}
}

func TestMaskValue_RevealTrailing(t *testing.T) {
	opts := MaskOptions{MaskString: "****", RevealChars: 3}
	result := MaskValue("supersecret", opts)
	if result != "****ret" {
		t.Errorf("expected '****ret', got %q", result)
	}
}

func TestMaskEntry_SensitiveKey(t *testing.T) {
	opts := DefaultMaskOptions()
	entry := Entry{Key: "DB_PASSWORD", Value: "hunter2"}
	result := MaskEntry(entry, opts)
	if result != "****" {
		t.Errorf("expected masked value, got %q", result)
	}
}

func TestMaskEntry_NonSensitiveKey(t *testing.T) {
	opts := DefaultMaskOptions()
	entry := Entry{Key: "APP_ENV", Value: "production"}
	result := MaskEntry(entry, opts)
	if result != "production" {
		t.Errorf("expected plain value, got %q", result)
	}
}

package crypto

import (
	"strings"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	if err := Init("test-key"); err != nil {
		t.Fatal(err)
	}

	plain := "super-secret-vapid-private-key"
	encrypted, err := Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(encrypted, prefix) {
		t.Fatalf("expected enc: prefix, got %q", encrypted)
	}
	if encrypted == plain {
		t.Fatal("encrypted value equals plaintext")
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if decrypted != plain {
		t.Fatalf("got %q, want %q", decrypted, plain)
	}
}

func TestDecryptPlaintextPassthrough(t *testing.T) {
	if err := Init("test-key"); err != nil {
		t.Fatal(err)
	}
	// Values written before encryption was enabled must pass through unchanged.
	v, err := Decrypt("legacy-plaintext-value")
	if err != nil {
		t.Fatal(err)
	}
	if v != "legacy-plaintext-value" {
		t.Fatalf("got %q", v)
	}
}

func TestEncryptedStringValueScan(t *testing.T) {
	if err := Init("test-key"); err != nil {
		t.Fatal(err)
	}

	original := EncryptedString("my-refresh-token")

	// Simulate DB write.
	dbVal, err := original.Value()
	if err != nil {
		t.Fatal(err)
	}
	stored, ok := dbVal.(string)
	if !ok {
		t.Fatalf("Value() returned non-string: %T", dbVal)
	}
	if !strings.HasPrefix(stored, prefix) {
		t.Fatalf("stored value missing enc: prefix: %q", stored)
	}

	// Simulate DB read.
	var result EncryptedString
	if err := result.Scan(stored); err != nil {
		t.Fatal(err)
	}
	if result != original {
		t.Fatalf("got %q, want %q", result, original)
	}
}

func TestNoKeyStoresPlaintext(t *testing.T) {
	block = nil // reset

	v, err := Encrypt("plaintext")
	if err != nil {
		t.Fatal(err)
	}
	if v != "plaintext" {
		t.Fatalf("expected passthrough when no key, got %q", v)
	}
}

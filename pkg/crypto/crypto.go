package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const prefix = "enc:"

var block cipher.Block

// Init derives a 32-byte AES key from key using SHA-256 and prepares the
// cipher block. Must be called before any Encrypt/Decrypt or EncryptedString
// DB operations. If key is empty, encryption is disabled and values are stored
// as plaintext (existing deployments can enable encryption and migrate gradually).
func Init(key string) error {
	if key == "" {
		return nil
	}
	sum := sha256.Sum256([]byte(key))
	b, err := aes.NewCipher(sum[:])
	if err != nil {
		return err
	}
	block = b
	return nil
}

// Encrypt returns an "enc:<base64>" string containing a random nonce prepended
// to the AES-256-GCM ciphertext. Returns plaintext unchanged when encryption
// is disabled (Init was called with an empty key).
func Encrypt(plaintext string) (string, error) {
	if block == nil || plaintext == "" {
		return plaintext, nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return prefix + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt reverses Encrypt. Values that do not start with "enc:" are returned
// as-is, which allows plaintext rows written before encryption was enabled to
// keep working without a bulk migration.
func Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, prefix) {
		return ciphertext, nil
	}
	if block == nil {
		return "", errors.New("crypto: ENCRYPTION_KEY not set but encrypted data found in database")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, prefix))
	if err != nil {
		return "", fmt.Errorf("crypto: base64 decode: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("crypto: ciphertext too short")
	}
	nonce, data := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", fmt.Errorf("crypto: decrypt: %w", err)
	}
	return string(plaintext), nil
}

// EncryptedString is a named string type that transparently encrypts on write
// and decrypts on read via the database/sql driver interfaces. Use it in place
// of string for sensitive GORM model fields.
type EncryptedString string

func (e EncryptedString) Value() (driver.Value, error) {
	if e == "" {
		return "", nil
	}
	return Encrypt(string(e))
}

func (e *EncryptedString) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case nil:
		*e = ""
		return nil
	default:
		return fmt.Errorf("crypto: cannot scan type %T into EncryptedString", value)
	}
	decrypted, err := Decrypt(s)
	if err != nil {
		return err
	}
	*e = EncryptedString(decrypted)
	return nil
}

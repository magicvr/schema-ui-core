package mail

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Secret-at-rest helpers for the R7 admin channel configuration
// (VP-017 / workspace-017 GOAL-008; Root D-007): secrets entered in the
// settings「邮件」tab are stored AES-256-GCM encrypted under a local master
// key and are NEVER returned by any read face ("secret 不入库明文可读" —
// VP-017 exit denominator). The master key never travels through the API.

// LoadOrCreateMasterKey resolves the 32-byte local master key:
//
//  1. envValue non-empty -> key = SHA-256(envValue) (operator-provided via
//     MAIL_CONFIG_MASTER_KEY; hashed so arbitrary-length passphrases work);
//  2. otherwise the 32 random bytes at keyPath are reused;
//  3. otherwise new random bytes are generated, written to keyPath with
//     0600 permissions, and returned.
//
// The auto-generated branch keeps the local two-process default path
// zero-config while still satisfying "not plaintext in the database".
func LoadOrCreateMasterKey(envValue, keyPath string) ([]byte, error) {
	if trimmed := strings.TrimSpace(envValue); trimmed != "" {
		sum := sha256.Sum256([]byte(trimmed))
		return sum[:], nil
	}
	if data, err := os.ReadFile(keyPath); err == nil && len(data) == masterKeyLen {
		return data, nil
	}
	key := make([]byte, masterKeyLen)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("mail: generate master key: %w", err)
	}
	if dir := filepath.Dir(keyPath); dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("mail: create master key directory: %w", err)
		}
	}
	if err := os.WriteFile(keyPath, key, 0o600); err != nil {
		return nil, fmt.Errorf("mail: persist master key: %w", err)
	}
	return key, nil
}

const masterKeyLen = 32

// EncryptSecret seals plaintext under key: base64(nonce || AES-GCM ciphertext).
// An empty plaintext yields an empty string (no secret stored), so untouched
// secret fields stay distinguishable from empty ones.
func EncryptSecret(key []byte, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("mail: encrypt secret: %w", err)
	}
	sealed := aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptSecret opens a value produced by EncryptSecret.
func DecryptSecret(key []byte, encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("mail: decode secret: %w", err)
	}
	aead, err := newAEAD(key)
	if err != nil {
		return "", err
	}
	if len(raw) < aead.NonceSize() {
		return "", fmt.Errorf("mail: decrypt secret: truncated payload")
	}
	nonce, ciphertext := raw[:aead.NonceSize()], raw[aead.NonceSize():]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("mail: decrypt secret: %w", err)
	}
	return string(plaintext), nil
}

func newAEAD(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("mail: init cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("mail: init gcm: %w", err)
	}
	return aead, nil
}

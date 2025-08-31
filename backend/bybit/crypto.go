package bybit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	encPrefix = "enc:v1:"
)

var (
	encKey []byte
)

func ensureEncKey() error {
	if encKey != nil {
		return nil
	}
	// 1) Try environment variable (dev)
	if key := os.Getenv("BYBIT_ENC_KEY"); key != "" {
		if k, err := base64.StdEncoding.DecodeString(key); err == nil {
			if len(k) != 32 {
				return fmt.Errorf("BYBIT_ENC_KEY decoded length must be 32 bytes; got %d", len(k))
			}
			encKey = k
			return nil
		}
		if len(key) != 32 {
			return fmt.Errorf("BYBIT_ENC_KEY length must be 32 bytes; got %d", len(key))
		}
		encKey = []byte(key)
		return nil
	}

	// 2) Try OS keyring (preferred in production installers)
	const service = "coin-control"
	const user = "enc-key"
	if stored, err := keyring.Get(service, user); err == nil && stored != "" {
		if k, err := base64.StdEncoding.DecodeString(stored); err == nil && len(k) == 32 {
			encKey = k
			return nil
		}
	}

	// 3) Generate and store new key in keyring
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return fmt.Errorf("failed generating enc key: %w", err)
	}
	b64 := base64.StdEncoding.EncodeToString(b)
	if err := keyring.Set(service, user, b64); err != nil {
		return fmt.Errorf("failed to store enc key in OS keyring: %w", err)
	}
	encKey = b
	return nil
}

func encryptString(plaintext string) (string, error) {
	if err := ensureEncKey(); err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	out := append(nonce, ciphertext...)
	return encPrefix + base64.StdEncoding.EncodeToString(out), nil
}

func decryptString(value string) (string, error) {
	// Backward compatibility: if value is not prefixed, assume plaintext
	if !strings.HasPrefix(value, encPrefix) {
		return value, nil
	}
	if err := ensureEncKey(); err != nil {
		return "", err
	}
	enc := strings.TrimPrefix(value, encPrefix)
	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}
	nonce := data[:gcm.NonceSize()]
	ct := data[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

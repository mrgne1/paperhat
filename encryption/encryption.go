package encryption

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var ErrInvalidKey error = errors.New("Invalid Key")

func Encrypt(value string) (string, string, error) {
	cipher := make([]byte, len(value))
	keyLength := len(value)
	keyLength += (10 - (keyLength % 10))
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		return "", "", err
	}

	for i := 0; i < len(value); i++ {
		cipher[i] = key[i] ^ value[i]
	}

	keyText := hex.EncodeToString(key)
	cipherText := hex.EncodeToString(cipher)
	return cipherText, keyText, nil
}

func Decrypt(cipherText string, keyText string) (string, error) {
	key, err := hex.DecodeString(keyText)
	if err != nil {
		return "", err
	}

	cipher, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(key) < len(cipher) {
		return "", ErrInvalidKey
	}

	clearText := make([]byte, len(cipher))
	for i := 0; i < len(cipher); i++ {
		clearText[i] = key[i] ^ cipher[i]
	}

	return string(clearText), nil
}

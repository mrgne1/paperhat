package encryption

import (
	"math/rand"
	"testing"
)

func TestEncryption(t *testing.T) {
	value := "A string to encrypt.!!?"
	_, _, err := Encrypt(value)

	if err != nil {
		t.Errorf("Unable to encrypt: %v", err)
	}
}

func TestEncryptionThenDecryption(t *testing.T) {
	value := "I have conquered the world with the power of truth"
	val, key, err := Encrypt(value)

	if err != nil {
		t.Errorf("Unable to encrypt: %v", err)
	}

	decrypted, err := Decrypt(val, key)

	if err != nil {
		t.Errorf("Unable to decrypt: %v", err)
	}

	checkEquality(t, decrypted, value)
}

func TestEncryptionSimulation(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		value := RandomString()
		val, key, err := Encrypt(value)

		if err != nil {
			t.Errorf("Unable to encrypt: %v", err)
		}

		decrypted, err := Decrypt(val, key)

		if err != nil {
			t.Errorf("Unable to decrypt: %v", err)
		}

		if decrypted != value {
			t.Errorf("Expect decrypted '%v' to equal original '%v'", decrypted, value)
		}
	}
}

func checkEquality(t *testing.T, decrypted, original string) {
	if len(decrypted) != len(original) {
		t.Errorf("Expected len(decrypted) '%v' to equal len(original) '%v'", len(decrypted), len(original))
	}

	if decrypted != original {

		for i := range decrypted {
			if decrypted[i] != original[i] {
				t.Errorf("Expected '%v' to equal '%v'\n", decrypted[i], original[i])
			}
		}
		// t.Errorf("Expect decrypted '%v' to equal original '%v'", decrypted, original)
	}
}

var charset string = "=+1!2@3#4$5%6^7&8*9(0)-_\tqQwWeErRtTyYuUiIoOpP\\|aAsSdDfFgGhHjJkKlL;:'\"zZxXcCvVbBnNmM,<.>/?[{]}"

func RandomString() string {
	length := rand.Intn(1000) + 1

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

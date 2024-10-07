package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func EncodeHashSha256(value string) string {
	h := hmac.New(sha256.New, []byte("KEY123!"))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func EncodeHashSha512(value string) string {
	h := hmac.New(sha512.New, []byte("KEY123!"))
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func encrypt(plaintext []byte, secretKey string) ([]byte, error) {
	tAES, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(tAES)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func decrypt(ciphertext []byte, secretKey string) ([]byte, error) {
	tAES, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(tAES)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strings"
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

func EncryptDecrypt(input, key string) string {
	kL := len(key)

	var tmp []string
	for i := 0; i < len(input); i++ {
		tmp = append(tmp, string(input[i]^key[i%kL]))
	}
	return strings.Join(tmp, "")
}

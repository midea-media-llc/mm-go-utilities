package utils

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
)

func Md5Hash(input string) (ouput string) {
	data := []byte(input)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func ByteArrayToBase64(bytes []byte) string {
	return b64.StdEncoding.EncodeToString(bytes)
}

func EncodeBase64(s string) string {
	data := b64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

func DecodeBase64(s string) (string, error) {
	data, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

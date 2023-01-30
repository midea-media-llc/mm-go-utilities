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

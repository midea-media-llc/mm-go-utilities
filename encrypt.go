package utils

import (
	"crypto/md5"
	b64 "encoding/base64"
	"fmt"
)

// Md5Hash function takes a string input and returns its MD5 hash as a hexadecimal string
func Md5Hash(input string) (ouput string) {
	// Convert input string to bytes
	data := []byte(input)

	// Hash the bytes using MD5 algorithm
	hash := md5.Sum(data)

	// Convert the hash to a hexadecimal string and return it
	return fmt.Sprintf("%x", hash)
}

// ByteArrayToBase64 function takes a byte slice and returns its base64 encoded string representation
func ByteArrayToBase64(bytes []byte) string {
	// Encode the byte slice using base64 and return the result as a string
	return b64.StdEncoding.EncodeToString(bytes)
}

// EncodeBase64 function takes a string input and returns its base64 encoded string representation
func EncodeBase64(s string) string {
	// Encode the input string as a byte slice using base64 and return the result as a string
	data := b64.StdEncoding.EncodeToString([]byte(s))
	return string(data)
}

// DecodeBase64 function takes a base64 encoded string input and returns the decoded string and any errors
func DecodeBase64(s string) (string, error) {
	// Decode the input string using base64
	data, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	// Convert the decoded bytes to a string and return it along with any errors
	return string(data), nil
}

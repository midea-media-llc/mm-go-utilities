package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash for a given password string
// with a cost of 2^14 iterations (recommended as of 2021).
// It takes the following argument:
// - password: the password string to hash
// The function returns the resulting hash string,
// or an error if the hashing operation fails.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash for the given password string with a cost of 2^14 iterations.
	// The cost controls the number of iterations and therefore the time it takes to hash the password.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Return the resulting hash string, or an error if the hashing operation fails.
	return string(bytes), err
}

// CheckPasswordHash checks whether a given password string matches a given hash string
// generated by the HashPassword function.
// It takes the following arguments:
// - password: the password string to check
// - hash: the hash string to compare against
// The function returns true if the password matches the hash,
// or false otherwise.
func CheckPasswordHash(password, hash string) bool {
	// Compare the given password string with the given hash string using bcrypt's CompareHashAndPassword function.
	// If the comparison succeeds, err will be nil; otherwise, it will contain an error message.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	// Return true if the comparison succeeds (i.e., err is nil), or false otherwise.
	return err == nil
}

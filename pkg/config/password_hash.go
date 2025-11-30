package config

import (
	"crypto/sha1"
	"fmt"
)

const salt = "4e7d2f82df296dggxwducvuwdvbcwy8egc98wyegHJKHVGkugvbdh"

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func CheckPasswordHash(password, hash string) bool {
	return GeneratePasswordHash(password) == hash
}
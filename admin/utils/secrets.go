package utils

import (
	"os"
	"fmt"
	"crypto/sha256"

	"github.com/TheAlpha16/isolet/admin/database"
)

func UpdateKey(key string) error {
	secret := database.GenerateRandom()
	err := os.Setenv(key, secret)
	if err != nil {
		return err
	}
	return nil
}

func Hash(secret string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(secret)))
}
package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/CyberLabs-Infosec/isolet/ripper/config"
)

func StringAddr(s string) *string {
	tempString := s
	return &tempString
}

func GetInstanceName(userid int, level int) string {
	return Hash(fmt.Sprintf("%d@%d:%s", userid, level, config.INSTANCE_NAME_SECRET))[0:16]
}

func Hash(secret string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(secret)))
}

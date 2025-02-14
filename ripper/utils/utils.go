package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/TheAlpha16/isolet/ripper/config"
)

func StringAddr(s string) *string {
	tempString := s
	return &tempString
}

func GetInstanceName(chall_id int, teamid int64) string {
	return Hash(fmt.Sprintf("%d@%d:%s", teamid, chall_id, config.INSTANCE_NAME_SECRET))[0:16]
}

func Hash(secret string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(secret)))
}

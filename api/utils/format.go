package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/TheAlpha16/isolet/api/config"
)

func GetInstanceName(chall_id int, teamid int64) string {
	return Hash(fmt.Sprintf("%d@%d:%s", teamid, chall_id, config.INSTANCE_NAME_SECRET))[0:16]
}

func GetHostName(chall_id int, teamid int64) string {
	return config.INSTANCE_HOSTNAME
}

func GetChallengeSubdomain(input string) string {
	subdomain := strings.ToLower(input)

	re := regexp.MustCompile(`[^a-z0-9-]`)
	subdomain = re.ReplaceAllString(subdomain, "-")

	subdomain = strings.Trim(subdomain, "-")

	if len(subdomain) > 63 {
		subdomain = subdomain[:63]
	}

	return subdomain
}

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func StringAddr(s string) *string {
	tempString := s
	return &tempString
}

func Int64Addr(i string) *int64 {
	tempInt, _ := strconv.ParseInt(i, 10, 64)
	return &tempInt
}

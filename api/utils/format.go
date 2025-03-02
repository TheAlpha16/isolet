package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/TheAlpha16/isolet/api/config"
)

func GetInstanceName(chall_id int, teamid int64, challenge_name string) string {
	return fmt.Sprintf("%s-%s", challenge_name, Hash(fmt.Sprintf("%d@%d:%s", teamid, chall_id, config.INSTANCE_NAME_SECRET))[0:16])
}

func GetHostName(subdomains []string) string {
	instanceHostName, _ := config.Get("INSTANCE_HOSTNAME")

	return strings.Join(subdomains, ".") + "." + instanceHostName
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

func Int32Addr(i string) *int32 {
	tempInt, _ := strconv.ParseInt(i, 10, 32)
	tempInt32 := int32(tempInt)
	return &tempInt32
}

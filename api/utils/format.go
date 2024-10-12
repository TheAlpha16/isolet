package utils

import (
	"fmt"
	"strconv"

	"github.com/TheAlpha16/isolet/api/config"
)

func GetInstanceName(chall_id int, teamid int64) string {
	return Hash(fmt.Sprintf("%d@%d:%s", teamid, chall_id, config.INSTANCE_NAME_SECRET))[0:16]
}

func GetHostName(userid int, level int) string {
	return config.INSTANCE_HOSTNAME
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

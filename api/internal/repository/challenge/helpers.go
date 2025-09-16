package challenge

import "fmt"

func GetChallengeCacheKey(id int64) string {
	return fmt.Sprintf("challenge:%d", id)
}

func GetAllChallengesCacheKey() string {
	return "challenges"
}

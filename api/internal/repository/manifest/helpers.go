package manifest

import "fmt"

func GetManifestCacheKey(challengeID int64) string {
	return fmt.Sprintf("challenge:%d:manifest", challengeID)
}

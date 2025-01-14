package database

import (
	"fmt"
	"crypto/rand"
	"encoding/hex"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/lib/pq"
)

func GenerateRandom() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
}

func isChallengeSolved(challengeID int64, solvedChalls pq.Int64Array) bool {
	for _, solved := range solvedChalls {
		if solved == challengeID {
			return true
		}
	}

	return false
}

func isRequirementMet(requirements pq.Int64Array, solvedChalls pq.Int64Array) bool {
	for _, requiredChall := range requirements {
		if !isChallengeSolved(requiredChall, solvedChalls) {
			return false
		}
	}

	return true
}

func GenerateChallengeEndpoint(method string, subdomain string, domain string, port int, username ...string) string {
	var connString string

	if subdomain != "" {
		subdomain = subdomain + "."
	}

	switch method {
		case "http":
			if port == 80 {
				connString = fmt.Sprintf("http://%s%s", subdomain, domain) 
			} else if port == 443 {
				connString = fmt.Sprintf("https://%s%s", subdomain, domain)
			} else {
				connString = fmt.Sprintf("http://%s%s:%d", subdomain, domain, port)
			}
		
		case "ssh":
			var user string

			if len(username) > 0 {
				user = username[0]
			} else {
				user = config.DEFAULT_USERNAME
			}

			if port == 22 {
				connString = fmt.Sprintf("ssh %s@%s%s", user, subdomain, domain)
			} else {
				connString = fmt.Sprintf("ssh %s@%s%s -p %d", user, subdomain, domain, port)
			}

		case "nc":
			connString = fmt.Sprintf("nc %s%s %d", subdomain, domain, port)
	}
	
	return connString
}
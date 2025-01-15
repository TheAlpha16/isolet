package database

import (
	"fmt"
	"crypto/rand"
	"encoding/hex"

	"github.com/TheAlpha16/isolet/api/config"
)

func GenerateRandom() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
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
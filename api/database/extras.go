package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/models"
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

func CleanSQLException(error_msg string) string {
	error_msg = strings.TrimPrefix(error_msg, "ERROR: ")

	prefixIndex := strings.Index(error_msg, " (SQLSTATE")
	if prefixIndex != -1 {
		error_msg = error_msg[:prefixIndex]
	}

	return error_msg
}

func GetConfigs(ctx context.Context) ([]models.Config, error) {
	db := DB.WithContext(ctx)
	
	var configs []models.Config
	if err := db.Table("config").Find(&configs).Error; err != nil {
		log.Printf("Failed to get configs: %v\n", err)
		return nil, err
	}

	return configs, nil
}

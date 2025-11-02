package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// StructTagsAsString collects struct field annotations up to a specified depth.
func StructTagsAsString(obj interface{}, tagKey string, depth int) string {
	var tags []string

	// Helper function to process struct fields recursively
	var collectTags func(reflect.Type, int)
	collectTags = func(t reflect.Type, currentDepth int) {
		if currentDepth > depth {
			return
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				// Recursively handle embedded structs
				collectTags(field.Type, currentDepth+1)
			} else {
				tag := field.Tag.Get(tagKey)
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("Provided type is not a struct")
	}

	collectTags(t, 0)
	return strings.Join(tags, ",")
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

// ComparePassword compares a hashed password with a plaintext password.
func ComparePassword(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

// RandomUUID generates a random UUID string.
func RandomUUID() string {
	return uuid.New().String()
}

// BuildLink constructs a URL with the given public URL, token, and path segments.
func BuildLink(publicURL, token string, pathSegments []string) (string, error) {
	url, err := url.Parse(publicURL)
	if err != nil {
		return "", err
	}

	url = url.JoinPath(pathSegments...)
	query := url.Query()
	query.Set(TokenQueryKey, token)
	url.RawQuery = query.Encode()

	if url.Scheme == "" {
		url.Scheme = "https"
	}

	return url.String(), nil
}

// IntOrNil returns a pointer to the integer or nil if the integer is zero.
func IntOrNil(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// Int64OrNil returns a pointer to the int64 or nil if the int64 is zero.
func Int64OrNil(i64 int64) *int64 {
	if i64 == 0 {
		return nil
	}
	return &i64
}

func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

// SlugifyForSubdomain generates a slug suitable for use as a subdomain.
func SlugifyForSubdomain(name string) string {
	re := regexp.MustCompile(RegexNonAlphanumeric)

	slug := strings.ToLower(name)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	// truncate to max 63 chars (DNS label limit)
	if len(slug) > 63 {
		slug = slug[:63]
	}

	slug = strings.TrimRight(slug, "-")

	return slug
}

// GenerateRandom generates a random 128-character hexadecimal string.
func GenerateRandom() string {
	buffer := make([]byte, 64)
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
}

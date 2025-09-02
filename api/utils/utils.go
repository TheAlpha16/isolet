package utils

import (
	"reflect"
	"strings"

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

func HashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

func RandomUUID() string {
	return uuid.New().String()
}

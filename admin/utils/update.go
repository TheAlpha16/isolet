package utils

import (
	"reflect"
	"strings"

	"github.com/TheAlpha16/isolet/admin/models"
	"github.com/lib/pq"
)


func UpdateChallenges(existingChallenge *models.Challenge, challengeMetaData *models.Challenge) *models.Challenge {
	if challengeMetaData.Name != "" {
		existingChallenge.Name = challengeMetaData.Name
	}

	if challengeMetaData.Author != "" {
		existingChallenge.Author = challengeMetaData.Author
	} 

	if challengeMetaData.Points > 0 {
		existingChallenge.Points = challengeMetaData.Points
	}

	if challengeMetaData.CategoryID > 0 {
		existingChallenge.CategoryID = challengeMetaData.CategoryID
	}

	if challengeMetaData.Prompt != "" {
		existingChallenge.Prompt = challengeMetaData.Prompt
	}

	if challengeMetaData.Type != "" {
		existingChallenge.Type = challengeMetaData.Type
	}

	if reflect.TypeOf(challengeMetaData.Visible).Kind() == reflect.Bool {
		existingChallenge.Visible = challengeMetaData.Visible
	}

	if challengeMetaData.Flag != "" {
		existingChallenge.Flag = challengeMetaData.Flag
	} 
	
	if challengeMetaData.Tags != nil {
		t := []string(challengeMetaData.Tags)
		dbTags := parseJsonStringArrays(t)

		existingChallenge.Tags = pq.StringArray(dbTags)
	}

	if challengeMetaData.Links != nil {
		l := []string(challengeMetaData.Links)
		dbLinks := parseJsonStringArrays(l)

		existingChallenge.Links = pq.StringArray(dbLinks)
	}

	return existingChallenge
}

func UpdateFiles (existingChallenge *models.Challenge, challengeMetaData *models.Challenge) *models.Challenge {
	if challengeMetaData.Files != nil {
		f := []string(challengeMetaData.Files)
		dbFiles := parseJsonStringArrays(f)

		existingChallenge.Links = pq.StringArray(dbFiles)
	}

	return existingChallenge
}

func parseJsonStringArrays(arr []string) []string {
	var dbArr []string

	trimmed := strings.Trim(arr[0], "[]")
	parts := strings.Split(trimmed, ",")
	for _, p := range parts {
		p := strings.Trim(p, " ")
		dbArr = append(dbArr, p)
	}

	return dbArr
}

// func parseJsonInt64Arrays(arr []string) []int64 {
// 	var dbArr []int64

// 	trimmed := strings.Trim(arr[0], "[]")
// 	parts := strings.Split(trimmed, ",")
// 	for _, p := range parts {
// 		p := strings.Trim(p, " ")
// 		pInt, _ := strconv.Atoi(p)
// 		dbArr = append(dbArr, int64(pInt))
// 	}

// 	return dbArr
// }


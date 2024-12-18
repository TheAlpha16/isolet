package database

import (
	"log"
	"time"
	"context"

	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
)

func ReadScores(c *fiber.Ctx, page int) (models.ScoreBoard, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	perPage := 10 
	offset := (page - 1) * perPage
	scores := make([]models.Score, 0)
	var totalTeams int64
	
	if err := db.Model(&models.Team{}).Count(&totalTeams).Error; err != nil {
		log.Println(err)
		return models.ScoreBoard{}, err
	}

	pageCount := int((totalTeams + int64(perPage) - 1) / int64(perPage))
	if pageCount < page {
		return models.ScoreBoard{
			PageCount: pageCount,
			Page:      page,
			Scores:    scores,
		}, nil
	}

	err := db.Table("teams").
		Select(`teams.teamid AS teamid, 
				teams.teamname AS teamname, 
				COALESCE(SUM(challenges.points), 0) - teams.cost AS score`).
		Joins("LEFT JOIN challenges ON challenges.chall_id = ANY(teams.solved)").
		Group("teams.teamid").
		Order("score DESC, teams.last_submission ASC").
		Limit(perPage).
		Offset(offset).
		Scan(&scores).Error

	if err != nil {
		log.Println(err)
		return models.ScoreBoard{}, err
	}

	for i := range scores {
		scores[i].Rank = offset + i + 1
	}

	return models.ScoreBoard{
		PageCount: pageCount,
		Page:      page,
		Scores:    scores,
	}, nil
}

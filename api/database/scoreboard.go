package database

import (
	"context"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/api/models"

	"github.com/gofiber/fiber/v2"
)

func ReadScores(c *fiber.Ctx, page int) (models.ScoreBoard, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	perPage := 50
	offset := (page - 1) * perPage
	var scores []models.ScoreBoardTeam
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

	if err := db.Model(&models.Team{}).Select("teamid, teamname, rank, score").Raw("SELECT teamid, teamname, rank, score FROM get_scoreboard(?, ?)", perPage, offset).Scan(&scores).Error; err != nil {
		log.Println(err)
		return models.ScoreBoard{}, err
	}

	return models.ScoreBoard{
		PageCount: pageCount,
		Page:      page,
		Scores:    scores,
	}, nil
}

func GetTeamScore(teamid int64) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var score int
	if err := db.Raw("SELECT calculate_score(?)", teamid).Scan(&score).Error; err != nil {
		log.Println(err)
		return 0, err
	}

	return score, nil
}

func GetScoreGraph(c *fiber.Ctx) ([]models.ScoreBoardTeam, error) {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	db := DB.WithContext(ctx)

	var scoreGraph []models.ScoreBoardTeam
	if err := db.Raw("SELECT teamid, teamname, rank, submissions FROM get_top_teams_submissions()").Scan(&scoreGraph).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if scoreGraph == nil {
		scoreGraph = []models.ScoreBoardTeam{}
	}

	return scoreGraph, nil
}

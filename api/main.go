package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
	"github.com/TheAlpha16/isolet/api/database"
	"github.com/TheAlpha16/isolet/api/logs"
	"github.com/TheAlpha16/isolet/api/router"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	logs.InitLogger()
	log.Println("API version: v2.0.0")

	for {
		if err := database.Connect(); err != nil {
			log.Println(err)
			log.Println("sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	// Generate new secrets if don't exist
	if err := utils.UpdateKey("SESSION_SECRET"); err != nil {
		log.Fatal(err)
	}
	if err := utils.UpdateKey("TOKEN_SECRET"); err != nil {
		log.Fatal(err)
	}
	if err := utils.UpdateKey("INSTANCE_NAME_SECRET"); err != nil {
		log.Fatal(err)
	}

	// Setup access logs
	accessLogFile, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogFile.Close()
	aw := io.MultiWriter(os.Stdout, accessLogFile)
	loggerConfig := logger.Config{
		Output: aw,
		Format: "${time} | ${status} | ${latency} | ${locals:clientIP} | ${method} | ${path} | ${error}\n",
	}

	app := fiber.New()
	app.Use(logger.New(loggerConfig))
	app.Use(recover.New())
	router.SetupRoutes(app)

	log.Fatal(app.Listen(config.APP_PORT))
}

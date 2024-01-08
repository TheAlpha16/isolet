package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/CyberLabs-Infosec/isolet/goapi/config"
	"github.com/CyberLabs-Infosec/isolet/goapi/database"
	"github.com/CyberLabs-Infosec/isolet/goapi/logs"
	"github.com/CyberLabs-Infosec/isolet/goapi/router"
	"github.com/CyberLabs-Infosec/isolet/goapi/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Initialize error logging
	logs.InitLogger()

	log.Println("API version: 0.9.7")
	// Connect to database
	for {
		if err := database.Connect(); err != nil {
			log.Println(err)
			log.Println("sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	// Create tables
	if err := database.CreateTables(); err != nil {
		log.Fatal(err)
	}

	// Initialize challenges
	if err := database.PopulateChalls(); err != nil {
		log.Fatal(err)
	}

	// Generate new secrets
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
	accessLogFile, err := os.OpenFile("./logs/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer accessLogFile.Close()
	aw := io.MultiWriter(os.Stdout, accessLogFile)
	loggerConfig := logger.Config{Output: aw}

	// Initialize *fiber.App
	app := fiber.New()
	app.Use(logger.New(loggerConfig)) // Add Logger middleware with config
	app.Use(recover.New())            // Prevent process exit due to Fatal()
	router.SetupRoutes(app)           // Setup routing

	log.Fatal(app.Listen(config.APP_PORT))
}

// TO-DO:
// 1. Use helm
// 2.
// 3. add resource contraints to pods
// 4.
// 5.
// 6.

// 99. Remove client access from inside cluster for pods
// 100. change app.Listen to app.ListenTLS

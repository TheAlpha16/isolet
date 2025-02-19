package main

import (
	"log"
	"time"

	"github.com/TheAlpha16/isolet/ripper/database"
	"github.com/TheAlpha16/isolet/ripper/instance"
)

func main() {
	log.Println("[LOG] Ripper version: v2.0.0")
	log.Println("[LOG] Connecting to DB...")

	for {
		if err := database.Connect(); err != nil {
			log.Printf("[ERROR] unable to connect to DB: %s\n", err.Error())
			log.Println("[LOG] sleep for 1 minute")
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	log.Printf("[LOG] DB connection established\n")

	instance.StartWatch()
	log.Println("[LOG] Ripper exiting...")
}

package main

import (
	"github.com/TheAlpha16/isolet/api/utils/logger"
)

func main() {
	logger.Init()
	defer logger.GetAppLogger().Sync()
}

package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var INSTANCE_NAME_SECRET = os.Getenv("INSTANCE_NAME_SECRET")
var INSTANCE_NAMESPACE = os.Getenv("INSTANCE_NAMESPACE")
var KUBECONFIG_FILE_PATH = os.Getenv("KUBECONFIG_FILE_PATH")
var WORKER_COUNT = os.Getenv("WORKER_COUNT")

func init() {
	if WORKER_COUNT == "" {
		WORKER_COUNT = "5"
	}
}
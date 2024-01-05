package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var CONCURRENT_INSTANCES, _ = strconv.Atoi(os.Getenv("CONCURRENT_INSTANCES"))
var INSTANCE_TIME, _ = strconv.Atoi(os.Getenv("INSTANCE_TIME"))
var INSTANCE_NAME_SECRET = os.Getenv("INSTANCE_NAME_SECRET")
var INSTANCE_NAMESPACE = os.Getenv("INSTANCE_NAMESPACE")
var KUBECONFIG_FILE_PATH = os.Getenv("KUBECONFIG_FILE_PATH")

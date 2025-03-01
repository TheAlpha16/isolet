package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var CONCURRENT_INSTANCES, _ = strconv.Atoi(os.Getenv("CONCURRENT_INSTANCES"))
var INSTANCE_NAME_SECRET = os.Getenv("INSTANCE_NAME_SECRET")
var DEFAULT_USERNAME = "hacker"
var INSTANCE_HOSTNAME = os.Getenv("INSTANCE_HOSTNAME")
var INSTANCE_TIME, _ = strconv.Atoi(os.Getenv("INSTANCE_TIME"))
var MAX_INSTANCE_TIME, _ = strconv.Atoi(os.Getenv("MAX_INSTANCE_TIME"))

var INSTANCE_NAMESPACE = os.Getenv("INSTANCE_NAMESPACE")
var IMAGE_REGISTRY = os.Getenv("IMAGE_REGISTRY")
var KUBECONFIG_FILE_PATH = os.Getenv("KUBECONFIG_FILE_PATH")

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

var INSTANCE_NAMESPACE = os.Getenv("INSTANCE_NAMESPACE")
var IMAGE_REGISTRY_PREFIX = os.Getenv("IMAGE_REGISTRY_PREFIX")
var KUBECONFIG_FILE_PATH = os.Getenv("KUBECONFIG_FILE_PATH")

var CPU_REQUEST = os.Getenv("CPU_REQUEST")
var MEMORY_REQUEST = os.Getenv("MEMORY_REQUEST")
var CPU_LIMIT = os.Getenv("CPU_LIMIT")
var MEMORY_LIMIT = os.Getenv("MEMORY_LIMIT")
var DISK_REQUEST = os.Getenv("DISK_REQUEST")
var DISK_LIMIT = os.Getenv("DISK_LIMIT")
var TERMINATION_PERIOD = os.Getenv("TERMINATION_PERIOD")
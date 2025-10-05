set shell := ["bash", "-cu"]
set dotenv-load

GO := "go"
GORUN_COMMAND := GO + " run"
GOTIDY_COMMAND := GO + " mod tidy"
REGISTRY := "docker.io/thealpha16"
API := "api"
UI := "ui"

# Default target
default: help

# Show available commands
help:
	@just --list

# --- Common commands ---

# Build the docker image
docker-build RESOURCE TAG:
	#!/usr/bin/env bash
	cd {{RESOURCE}}
	echo "[#] building docker image for {{RESOURCE}} with tag {{TAG}}"
	docker build -t {{REGISTRY}}/isolet-{{RESOURCE}}:{{TAG}} .
	docker build -t {{REGISTRY}}/isolet-{{RESOURCE}}:latest .
	echo "[#] built docker image for {{RESOURCE}} with tag {{TAG}}"

# Push the docker image to the registry
docker-push RESOURCE TAG:
	#!/usr/bin/env bash
	if ! docker image inspect {{REGISTRY}}/isolet-{{RESOURCE}}:{{TAG}} > /dev/null 2>&1; then
		echo "[!] docker image for {{RESOURCE}} with tag {{TAG}} not found locally. Building it first..."
		just docker-build {{RESOURCE}} {{TAG}}
	fi
	echo "[#] pushing docker image for {{RESOURCE}} with tag {{TAG}} to registry"
	docker push {{REGISTRY}}/isolet-{{RESOURCE}}:{{TAG}}
	docker push {{REGISTRY}}/isolet-{{RESOURCE}}:latest
	echo "[#] pushed docker image for {{RESOURCE}} with tag {{TAG}} to registry"

# Bump version (usage: just bump RESOURCE patch | minor | major)
bump RESOURCE LEVEL:
	#!/usr/bin/env bash
	OLD=$(cat {{RESOURCE}}/VERSION)
	NEW=$(semver --increment {{LEVEL}} $OLD)
	echo $NEW > {{RESOURCE}}/VERSION
	echo "Bumped {{RESOURCE}} version: $OLD → $NEW"
	git add {{RESOURCE}}/VERSION
	git commit -m "chore({{RESOURCE}}): bump version to $NEW"

# --- API commands ---

# Run the API service
api-run:
	dotenv -f .env -- dotenv -f {{API}}/.env -- \
		sh -c 'cd {{API}} && {{GORUN_COMMAND}} main.go'

# Tidy API Go modules
api-tidy:
	cd {{API}} && {{GOTIDY_COMMAND}}

# Show current API version
api-version:
	@cat {{API}}/VERSION

# --- UI commands ---

# Generate UI client from OpenAPI spec
ui-generate:
	dotenv -f .env -- \
		npx openapi-typescript-codegen --input {{API}}/openapi.yaml --output {{UI}}/api
	@echo "Generated UI client from OpenAPI spec."

# Run the UI development server
ui-run:
	sh -c 'cd {{UI}} && npm install && npm run dev'

# Show current UI version
ui-version:
	@cat {{UI}}/VERSION

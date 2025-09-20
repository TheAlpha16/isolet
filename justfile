set shell := ["bash", "-cu"]
set dotenv-load

GO := "go"
GORUN_COMMAND := GO + " run"
GOTIDY_COMMAND := GO + " mod tidy"
API := "api"
UI := "ui"

# Default target
default: help

# Show available commands
help:
	@just --list

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

# Bump version (usage: just api-bump patch | minor | major)
api-bump LEVEL:
	#!/usr/bin/env bash
	OLD=$(cat {{API}}/VERSION)
	NEW=$(semver --increment {{LEVEL}} $OLD)
	echo $NEW > {{API}}/VERSION
	echo "Bumped API version: $OLD → $NEW"
	git add {{API}}/VERSION
	git commit -m "chore(api): bump version to $NEW"

# --- UI commands ---

# Generate UI client from OpenAPI spec
ui-generate:
	dotenv -f .env -- \
		npx openapi-typescript-codegen --input {{API}}/openapi.yaml --output {{UI}}/api
	@echo "Generated UI client from OpenAPI spec."

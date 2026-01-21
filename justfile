set shell := ["bash", "-cu"]
set dotenv-load

GO := "go"
GORUN_COMMAND := GO + " run"
GOTIDY_COMMAND := GO + " mod tidy"
REGISTRY := "docker.io/thealpha16"
ORACLE := "oracle"
UI := "ui"
TIDE := "tide"
PROXY := "proxy"
SOCKY := "socky"
HERALD := "herald"

# Default target
default: help

# Show available commands
help:
	@just --list

# --- Common commands ---

# Build the docker image
docker-build RESOURCE TAG="":
	#!/usr/bin/env bash
	if [ -z "{{TAG}}" ]; then
		TAG=$(cat {{RESOURCE}}/VERSION)
	else
		TAG="{{TAG}}"
	fi
	echo "[#] ensuring buildx builder exists..."
	docker buildx create --name isolet-builder --use 2>/dev/null || docker buildx use isolet-builder
	echo "[#] building multi-platform docker image for {{RESOURCE}} with tag $TAG"
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t {{REGISTRY}}/isolet-{{RESOURCE}}:$TAG \
		-t {{REGISTRY}}/isolet-{{RESOURCE}}:latest \
		-f {{RESOURCE}}/Dockerfile .
	echo "[#] built docker image for {{RESOURCE}} with tag $TAG"

# Push the docker image to the registry
docker-push RESOURCE TAG="":
	#!/usr/bin/env bash
	cd {{RESOURCE}}
	if [ -z "{{TAG}}" ]; then
		TAG=$(cat VERSION)
	else
		TAG="{{TAG}}"
	fi
	if ! docker image inspect {{REGISTRY}}/isolet-{{RESOURCE}}:$TAG > /dev/null 2>&1; then
		echo "[!] docker image for {{RESOURCE}} with tag $TAG not found locally. Building it first..."
		just docker-build {{RESOURCE}} $TAG
	fi
	echo "[#] pushing docker image for {{RESOURCE}} with tag $TAG to registry"
	docker push {{REGISTRY}}/isolet-{{RESOURCE}}:$TAG
	docker push {{REGISTRY}}/isolet-{{RESOURCE}}:latest
	echo "[#] pushed docker image for {{RESOURCE}} with tag $TAG to registry"

# Bump version (usage: just bump RESOURCE patch | minor | major)
bump RESOURCE LEVEL:
	#!/usr/bin/env bash
	OLD=$(cat {{RESOURCE}}/VERSION)
	NEW=$(semver --increment {{LEVEL}} $OLD)
	echo $NEW > {{RESOURCE}}/VERSION
	echo "Bumped {{RESOURCE}} version: $OLD → $NEW"
	git add {{RESOURCE}}/VERSION
	git commit -m "chore({{RESOURCE}}): bump version to $NEW"

# Build all services
build-all:
	just docker-build {{ORACLE}}
	just docker-build {{UI}}
	just docker-build {{TIDE}}
	just docker-build {{PROXY}}
	just docker-build {{SOCKY}}
	just docker-build {{HERALD}}
# Push all services
push-all:
	just docker-push {{ORACLE}}
	just docker-push {{UI}}
	just docker-push {{TIDE}}
	just docker-push {{PROXY}}
	just docker-push {{SOCKY}}
	just docker-push {{HERALD}}
# --- API commands ---

# Run the API service
oracle-run:
	dotenv -f .env -- dotenv -f {{ORACLE}}/.env -- \
		sh -c 'cd {{ORACLE}} && {{GORUN_COMMAND}} main.go'

# Tidy API Go modules
oracle-tidy:
	cd {{ORACLE}} && {{GOTIDY_COMMAND}}

# Show current API version
oracle-version:
	@cat {{ORACLE}}/VERSION

# --- Tide Controller commands ---

# Run the Tide controller locally (for development)
tide-run:
	cd {{TIDE}} && make run

# Install Tide CRDs into the cluster
tide-install-crds:
	cd {{TIDE}} && make install

# Uninstall Tide CRDs from the cluster
tide-uninstall-crds:
	cd {{TIDE}} && make uninstall

# Deploy Tide controller to the cluster
tide-deploy:
	cd {{TIDE}} && make deploy

# Undeploy Tide controller from the cluster
tide-undeploy:
	cd {{TIDE}} && make undeploy

# Run Tide tests
tide-test:
	cd {{TIDE}} && make test

# Tidy Tide Go modules
tide-tidy:
	cd {{TIDE}} && {{GOTIDY_COMMAND}}

# Show current Tide version
tide-version:
	@cat {{TIDE}}/VERSION

# --- UI commands ---

# Generate UI client from OpenAPI spec
ui-generate:
	dotenv -f .env -- \
		npx openoracle-typescript-codegen --input {{ORACLE}}/openapi.yaml --output {{UI}}/api
	@echo "Generated UI client from OpenAPI spec."

# Run the UI development server
ui-run:
	sh -c 'cd {{UI}} && npm install && npm run dev'

# Build the UI for production
ui-build:
	sh -c 'cd {{UI}} && npm install && npm run build'

# Show current UI version
ui-version:
	@cat {{UI}}/VERSION

# --- Proxy commands ---

# Show current Proxy version
proxy-version:
	@cat {{PROXY}}/VERSION

# --- Socky commands ---

# Run the Socky service locally
socky-run:
	sh -c 'cd {{SOCKY}} && npm install && npm run dev'

# Build the Socky service
socky-build:
	sh -c 'cd {{SOCKY}} && npm install && npm run build'

# Show current Socky version
socky-version:
	@cat {{SOCKY}}/VERSION

# --- Herald commands ---

# Run the Herald service
herald-run:
	dotenv -f .env -- dotenv -f {{HERALD}}/.env -- \
		sh -c 'cd {{HERALD}} && {{GORUN_COMMAND}} main.go'

# Tidy Herald Go modules
herald-tidy:
	cd {{HERALD}} && {{GOTIDY_COMMAND}}

# Show current Herald version
herald-version:
	@cat {{HERALD}}/VERSION

# --- Kubernetes/Helm commands ---

# Install/upgrade the Helm chart
helm-install:
	helm upgrade --install isolet ./charts -f ./charts/values.yaml

# Dry-run the Helm chart
helm-dry-run:
	helm upgrade --install isolet ./charts -f ./charts/values.yaml --dry-run --debug

# Uninstall the Helm chart
helm-uninstall:
	helm uninstall isolet

# Template the Helm chart (useful for debugging)
helm-template:
	helm template isolet ./charts -f ./charts/values.yaml

# Apply Instance CRD
apply-crds:
	kubectl apply -f charts/crds/instances.yaml
	kubectl apply -f charts/crds/traefik-crds.yaml

# Delete Instance CRD
delete-crds:
	kubectl delete -f charts/crds/instances.yaml
	kubectl delete -f charts/crds/traefik-crds.yaml

# --- Development commands ---

# Clean up generated files and build artifacts
clean:
	rm -rf {{ORACLE}}/tmp
	rm -rf {{UI}}/.next
	rm -rf {{UI}}/node_modules
	rm -rf {{SOCKY}}/node_modules
	rm -rf {{TIDE}}/bin

# Show all versions
versions:
	@echo "API:   $(cat {{ORACLE}}/VERSION)"
	@echo "UI:    $(cat {{UI}}/VERSION)"
	@echo "Tide:  $(cat {{TIDE}}/VERSION)"
	@echo "Proxy: $(cat {{PROXY}}/VERSION)"
	@echo "Socky: $(cat {{SOCKY}}/VERSION)"
	@echo "Herald: $(cat {{HERALD}}/VERSION)"

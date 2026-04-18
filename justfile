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
PULSE := "pulse"
HERALD := "herald"

# Default target
default: help

# Show available commands
help:
	@just --list

# --- Common commands ---

# Set up a multi-platform buildx builder.
setup-builder CERT="":
	#!/usr/bin/env bash
	set -euo pipefail
	if ! docker buildx inspect isolet-builder > /dev/null 2>&1; then
		docker buildx create --name isolet-builder --driver docker-container --bootstrap --use
		echo "[#] created buildx builder 'isolet-builder'"
	else
		docker buildx inspect --bootstrap isolet-builder > /dev/null
		docker buildx use isolet-builder
	fi
	if [ -n "{{CERT}}" ] && [ -f "{{CERT}}" ]; then
		CONTAINER="buildx_buildkit_isolet-builder0"
		echo "[#] installing CA cert into buildx builder"
		docker cp "{{CERT}}" "$CONTAINER:/tmp/extra.crt"
		docker exec "$CONTAINER" sh -c \
			"cat /tmp/extra.crt >> /etc/ssl/certs/ca-certificates.crt"
		docker restart "$CONTAINER" > /dev/null
		sleep 2
		echo "[#] builder ready with CA cert"
	fi

# Build and push the docker image (linux/amd64 + linux/arm64).
docker-build RESOURCE TAG="" CERT="":
	#!/usr/bin/env bash
	set -euo pipefail
	just setup-builder {{CERT}}
	cd {{RESOURCE}}
	if [ -z "{{TAG}}" ]; then
		TAG=$(cat VERSION)
	else
		TAG="{{TAG}}"
	fi
	BUILD_ARGS=()
	if [ -n "{{CERT}}" ] && [ -f "{{CERT}}" ]; then
		BUILD_ARGS+=(--build-arg "EXTRA_CA_CERT=$(cat \"{{CERT}}\")")
	fi
	echo "[#] building and pushing docker image for {{RESOURCE}} with tag $TAG (linux/amd64,linux/arm64)"
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		"${BUILD_ARGS[@]}" \
		-t {{REGISTRY}}/isolet-{{RESOURCE}}:$TAG \
		-t {{REGISTRY}}/isolet-{{RESOURCE}}:latest \
		--push \
		.
	echo "[#] pushed {{RESOURCE}}:$TAG"

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
build-all CERT="":
	just docker-build {{ORACLE}} "" {{CERT}}
	just docker-build {{UI}} "" {{CERT}}
	just docker-build {{TIDE}} "" {{CERT}}
	just docker-build {{PROXY}} "" {{CERT}}
	just docker-build {{PULSE}} "" {{CERT}}
	just docker-build {{HERALD}} "" {{CERT}}

# --- API commands ---

# Run the API service
oracle-run:
	cd {{ORACLE}} && {{GORUN_COMMAND}} main.go

# Tidy API Go modules
oracle-tidy:
	cd {{ORACLE}} && {{GOTIDY_COMMAND}}

# Show current API version
oracle-version:
	@cat {{ORACLE}}/VERSION

# Start consumer for processing
oracle-consumer:
	cd {{ORACLE}} && IDENTITY=consumer METRICS_PORT=6868 {{GORUN_COMMAND}} main.go

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
	npx openapi-typescript-codegen --input {{ORACLE}}/openapi.yaml --output {{UI}}/api
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

# Restart nginx proxy with updated configuration
proxy-restart:
	@bash dev/proxy/proxy-restart.sh

# Show current Proxy version
proxy-version:
	@cat {{PROXY}}/VERSION

# --- Pulse commands ---

# Run the Pulse service locally
pulse-run:
	sh -c 'cd {{PULSE}} && mix run'

# Build the Pulse service
pulse-build:
	sh -c 'cd {{PULSE}} && mix deps.get && mix compile'

# Show current Pulse version
pulse-version:
	@cat {{PULSE}}/VERSION

# --- Herald commands ---

# Run the Herald service with k8s source
herald-k8s:
	cd {{HERALD}} && IDENTITY=k8s_source {{GORUN_COMMAND}} main.go

# Run the Herald service with postgres source
herald-postgres:
	cd {{HERALD}} && IDENTITY=postgres_source {{GORUN_COMMAND}} main.go

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
	rm -rf {{PULSE}}/node_modules
	rm -rf {{TIDE}}/bin

# Show all versions
versions:
	@echo "API:   $(cat {{ORACLE}}/VERSION)"
	@echo "UI:    $(cat {{UI}}/VERSION)"
	@echo "Tide:  $(cat {{TIDE}}/VERSION)"
	@echo "Proxy: $(cat {{PROXY}}/VERSION)"
	@echo "Pulse: $(cat {{PULSE}}/VERSION)"
	@echo "Herald: $(cat {{HERALD}}/VERSION)"

-include .env
export

GIN_MODE=test
BUILD_TIME?=$(shell TZ=${TZ} date '+%Y-%m-%d %H:%M:%S')

LDFLAGS=$(shell echo \
	"-X 'osoc/pkg/application.buildVersionTime=${BUILD_TIME}'" \
)

DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-27s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Format golang files with goimports
	find . -name \*.go -not -path \*/wire_gen.go -exec goimports -w {} \;

.PHONY: finalcheck
finalcheck: wire fmt mod lint test-short swagger-gen ## Make a final complex check before the commit

.PHONY: run
run: ## Run project for local
	go run -race -ldflags="${LDFLAGS}" ./cmd/${APP_NAME}/.
	#go run -ldflags="${LDFLAGS}" ./cmd/${APP_NAME}/. 2> trace.out

.PHONY: debug
debug: ## Run all container without app
	docker-compose up -d --scale app=0
	make run

.PHONY: watch
watch: ## Run in live-reload mode
	make stop
	docker-compose up

.PHONY: stop
stop: ## Remove containers but keep volumes
	docker-compose down --remove-orphans

.PHONY: clear
clear: ### Remove containers and volumes
	docker-compose down --remove-orphans --volumes

.PHONY: rebuild
rebuild: ## Rebuild by Docker Compose
	make stop
	docker-compose build --no-cache


.PHONY: compile
compile: ## Make binary and docs
	go build -ldflags="${LDFLAGS}" -o bin/${APP_NAME} cmd/${APP_NAME}/main.go cmd/${APP_NAME}/wire_gen.go


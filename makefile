# Project variables
BINARY_NAME:=appointment-service
BUILD_DIR:=build
MAIN_PATH:=./cmd/api
VERSION:=$(shell cat VERSION)
COMMIT_SHA:=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME:=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_BIN:=$(shell if [ -n "$(shell go env GOBIN)" ]; then echo "$(shell go env GOBIN)"; else echo "$(shell go env GOPATH)/bin"; fi)

# Declare all PHONY targets
.PHONY: all build run clean test test-verbose \
        coverage-generation coverage-report-functions coverage-report-packages \
        coverage-report-total coverage-all coverage-browser \
        dependencies install-tools lint dev vet fmt quality \
        run-dev-mem run-dev-postgres check-air check-docker \
        tidy start-telemetry stop-telemetry follow-telemetry

#---------------------------------------------------------
# Environment loading function
#---------------------------------------------------------
define load_env
    $(eval ENV_FILE := .environments/$(1))
    $(if $(wildcard $(ENV_FILE)), \
        $(eval include $(ENV_FILE)) \
        $(eval export $(shell sed 's/=.*//' $(ENV_FILE))) \
        $(info Loading environment from $(ENV_FILE):) \
        $(foreach v,$(shell sed 's/=.*//' $(ENV_FILE)),$(info $(v)=$($(v)))) \
        $(eval ENV_VARS := $(shell cat $(ENV_FILE) | tr '\n' ' ')), \
        $(error Environment file $(ENV_FILE) not found))
endef

#---------------------------------------------------------
# Development and build targets
# Note: The ldflags are used to inject the version, 
# commit and build time.  These values are set at can be
# used in logging and other places in the application.
#---------------------------------------------------------
build: tidy fmt vet
	@echo "Building the Go service..."
	go build -ldflags "\
        -X appointment-service/internal/version.Version=${VERSION} \
        -X appointment-service/internal/version.Commit=${COMMIT_SHA} \
        -X appointment-service/internal/version.BuildTime=${BUILD_TIME}" \
        -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api
	@echo "Build complete. Output: $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@echo "Starting the Go service..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
	go clean -testcache
	rm -f coverage.out
	rm -f coverage.html
	rm -f .env.tmp 

#---------------------------------------------------------
# Code quality targets
#---------------------------------------------------------
tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

quality: fmt vet lint test

#---------------------------------------------------------
# Test targets
#---------------------------------------------------------
test:
	go test -race ./...

test-verbose:
	go test -race -v ./...

#---------------------------------------------------------
# Coverage targets
#---------------------------------------------------------
coverage-generation:
	@echo "================================="
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./... >/dev/null 
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Done! "
	@echo "  raw results --> ./coverage.out"
	@echo "  html report --> ./coverage.html"

coverage-report-functions: coverage-generation
	@echo
	@echo "================================="
	@echo "  Results Per Function: "
	@echo "================================="
	@go tool cover -func=coverage.out

coverage-report-packages: coverage-generation
	@echo
	@echo "================================="
	@echo " Results Per package: "
	@echo "================================="
	@go test -cover ./...

coverage-report-total: coverage-generation
	@echo
	@echo "================================="
	@go tool cover -func=coverage.out | grep total: | awk '{print " TOTAL TEST COVERGAGE:  " $$3}'
	@echo "================================="

coverage-all: coverage-generation coverage-report-functions coverage-report-packages coverage-report-total

coverage-browser: coverage-generation
	open coverage.html

#---------------------------------------------------------
# Development environment targets
#---------------------------------------------------------
check-air:
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Error: air not found in path. Please ensure $(GO_BIN) is in your PATH"; \
		exit 1; \
	fi

run-dev-mem: check-air build
	$(call load_env,development.memory)
	@echo "Starting the Go service in development (in-memory storage) mode with live reload..."
	env $(ENV_VARS) air -c .air.toml

run-dev-postgres: check-air build
	$(call load_env,development.postgres)
	@echo "Starting the Go service in development (postgres storage) mode with live reload..."
	env $(ENV_VARS) air -c .air.toml

run-dev-sqlite3: check-air build
	$(call load_env,development.sqlite3)
	@echo "Starting the Go service in development (sqlite3 storage) mode with live reload..."
	env $(ENV_VARS) air -c .air.toml

check-docker:
	@if ! command -v docker > /dev/null 2>&1; then \
		echo "Docker is not installed. Please install."; \
		exit 1; \
	fi

# TODO: add support for postgres via docker
# start-dev-postgres: check-docker
# 	@echo "Starting development postgres database..."

#---------------------------------------------------------
# DB migration targets
#---------------------------------------------------------
DB_PATH ?= data/appointments.db
MIGRATIONS_DIR = migrations

# Apply all up migrations
migrate-up:
	migrate -database "sqlite3://$(DB_PATH)" -path $(MIGRATIONS_DIR) up

# Rollback all migrations
migrate-down:
	echo "y" | migrate -database "sqlite3://data/appointments.db" -path migrations down

# Clean database
clean-db:
	rm -f $(DB_PATH)

# Rebuild database from scratch
rebuild-db: clean-db migrate-up


#---------------------------------------------------------
# Tool installation
#---------------------------------------------------------
install-tools:
	@echo "Installing required tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed to $(GO_BIN)"
	@echo "Ensure $(GO_BIN) is in your PATH"
